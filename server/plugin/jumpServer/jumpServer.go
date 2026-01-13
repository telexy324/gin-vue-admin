package jumpServer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/google/uuid"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type Session struct {
	ID           string
	User         string
	Target       string
	StartedAt    time.Time
	LastActiveAt atomic.Int64

	Ctx    context.Context
	Cancel context.CancelFunc
}

type activityWriter struct {
	rw      io.ReadWriter
	session *Session
}

func (a *activityWriter) Write(p []byte) (int, error) {
	n, err := a.rw.Write(p)
	a.session.LastActiveAt.Store(time.Now().Unix())
	if err != nil {
		a.session.Cancel()
	}
	return n, err
}

func (a *activityWriter) Read(p []byte) (int, error) {
	n, err := a.rw.Read(p)
	a.session.LastActiveAt.Store(time.Now().Unix())
	if err != nil {
		a.session.Cancel()
	}
	return n, err
}

type ProxySFTPHandler struct {
	client *sftp.Client
}

/******** FileReader ********/
func (h *ProxySFTPHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	return h.client.Open(r.Filepath)
}

/******** FileWriter ********/
func (h *ProxySFTPHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	return h.client.Create(r.Filepath)
}

/******** FileCmder ********/
func (h *ProxySFTPHandler) Filecmd(r *sftp.Request) error {
	switch r.Method {
	case "Remove":
		return h.client.Remove(r.Filepath)
	case "Mkdir":
		return h.client.Mkdir(r.Filepath)
	case "Rename":
		return h.client.Rename(r.Filepath, r.Target)
	case "Rmdir":
		return h.client.RemoveDirectory(r.Filepath)
	case "Setstat":
		return nil
	default:
		return fmt.Errorf("unsupported cmd: %s", r.Method)
	}
}

/******** FileLister ********/
func (h *ProxySFTPHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
	files, err := h.client.ReadDir(r.Filepath)
	if err != nil {
		return nil, err
	}
	return &fileInfoLister{files: files}, nil
}

type fileInfoLister struct {
	files []os.FileInfo
}

func (l *fileInfoLister) ListAt(dst []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(l.files)) {
		return 0, io.EOF
	}
	n := copy(dst, l.files[offset:])
	return n, nil
}

const (
	idleTimeoutDefault    = 10 * time.Minute
	maxSessionTimeDefault = 2 * time.Hour
)

var targetPool sync.Map // map[string]*ssh.Client
var sftpPool sync.Map   // key: targetName, value: *sftp.Client

func Init() {
	// 1. SSH Server 配置
	homePath, err := os.UserHomeDir()
	if err != nil {
		global.GVA_LOG.Fatal("get home path fail", zap.Any("jump server", err))
	}
	privateBytes, err := os.ReadFile(path.Join(homePath, ".ssh", "id_rsa"))
	if err != nil {
		global.GVA_LOG.Fatal("get private key file fail", zap.Any("jump server", err))
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		global.GVA_LOG.Fatal("parse private key fail", zap.Any("jump server", err))
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true, // 简化示例（生产请做认证）
	}
	config.AddHostKey(private)

	// 2. 监听 22（或其他端口）
	port := 22
	if global.GVA_CONFIG.JumpServer.Port > 0 {
		port = global.GVA_CONFIG.JumpServer.Port
	}
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		global.GVA_LOG.Fatal("server listen fail", zap.Any("jump server", err))
	}
	global.GVA_LOG.Info("Jump server listening on ", zap.String("addr", addr))

	for {
		conn, _ := listener.Accept()
		go handleConn(conn, config)
	}
}

func handleConn(nConn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		global.GVA_LOG.Error("handshake failed: ", zap.Any("jump server", err))
		return
	}
	defer sshConn.Close()

	targetName := sshConn.User()
	sess := newSession("", targetName)

	// 3. 连接目标服务器
	targetClient, err := connectTarget(targetName)
	if err != nil {
		global.GVA_LOG.Error("target connect failed: ", zap.Any("jump server", err))
		return
	}
	defer targetClient.Close()

	go ssh.DiscardRequests(reqs)

	for ch := range chans {
		if ch.ChannelType() != "session" {
			ch.Reject(ssh.UnknownChannelType, "")
			continue
		}

		srcChannel, srcRequests, _ := ch.Accept()
		dstSession, err := targetClient.NewSession()
		if err != nil {
			global.GVA_LOG.Error("create session failed: ", zap.Any("jump server", err))
			srcChannel.Close()
			continue
		}

		go func() {
			for req := range srcRequests {
				switch req.Type {

				case "pty-req":
					dstSession.RequestPty("xterm-256color", 40, 120, ssh.TerminalModes{})
					req.Reply(true, nil)

				case "shell":
					//dstSession.Stdin = srcChannel
					//dstSession.Stdout = srcChannel
					//dstSession.Stderr = srcChannel
					dstSession.Stdout = &activityWriter{rw: srcChannel, session: sess}
					dstSession.Stderr = &activityWriter{rw: srcChannel, session: sess}
					dstSession.Stdin = &activityWriter{rw: srcChannel, session: sess}

					dstSession.Shell()
					go func() {
						err := dstSession.Wait()
						log.Println("target session exited:", err)
						sess.Cancel()
					}()
					req.Reply(true, nil)

				case "subsystem":
					if string(req.Payload[4:]) == "sftp" {
						req.Reply(true, nil)
						startSFTPServer(srcChannel, targetClient, targetName)
						return
					}
					req.Reply(false, nil)

				}
			}
		}()

		go monitorSession(sess)

		<-sess.Ctx.Done()
		srcChannel.Close()
		dstSession.Close()
	}
}

func connectTarget(targetName string) (*ssh.Client, error) {
	// 🍀 尝试复用
	if val, ok := targetPool.Load(targetName); ok {
		client := val.(*ssh.Client)
		// 测试是否还活着
		_, _, err := client.SendRequest("keepalive@golang.org", true, nil)
		if err == nil {
			global.GVA_LOG.Info("[reuse] reuse ssh client for ", zap.String("target", targetName))
			return client, nil
		}

		// 已失效，关闭并删掉
		client.Close()
		targetPool.Delete(targetName)
	}

	// 🍀 建立新连接
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	key, err := os.ReadFile(path.Join(homePath, ".ssh", "id_rsa"))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", targetName, 1122)
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}

	// 🍀 放入复用池
	targetPool.Store(targetName, client)
	global.GVA_LOG.Info("[new] new ssh client connected: ", zap.String("target", targetName))

	return client, nil
}

func newSession(user, target string) *Session {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Session{
		ID:        uuid.NewString(),
		User:      user,
		Target:    target,
		StartedAt: time.Now(),
		Ctx:       ctx,
		Cancel:    cancel,
	}
	s.LastActiveAt.Store(time.Now().Unix())
	return s
}

func monitorSession(sess *Session) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	var maxSessionTime, idleTimeout = maxSessionTimeDefault, idleTimeoutDefault
	if global.GVA_CONFIG.JumpServer.MaxSessionTime > 0 {
		maxSessionTime = time.Duration(global.GVA_CONFIG.JumpServer.MaxSessionTime) * time.Hour
	}
	if global.GVA_CONFIG.JumpServer.IdleTimeout > 0 {
		idleTimeout = time.Duration(global.GVA_CONFIG.JumpServer.IdleTimeout) * time.Minute
	}

	for {
		select {
		case <-sess.Ctx.Done():
			return

		case <-ticker.C:
			now := time.Now()

			// 最大会话时长
			if now.Sub(sess.StartedAt) > maxSessionTime {
				global.GVA_LOG.Info("session max time reached: ", zap.String("jump server", sess.ID))
				sess.Cancel()
				return
			}

			// 空闲超时
			last := time.Unix(sess.LastActiveAt.Load(), 0)
			if now.Sub(last) > idleTimeout {
				global.GVA_LOG.Info("session idle timeout: ", zap.String("jump server", sess.ID))
				sess.Cancel()
				return
			}
		}
	}
}

func getSftpClient(targetName string, targetClient *ssh.Client) (*sftp.Client, error) {
	if val, ok := sftpPool.Load(targetName); ok {
		cli := val.(*sftp.Client)
		// 测试是否存活
		_, err := cli.ReadDir("/")
		if err == nil {
			return cli, nil
		}

		cli.Close()
		sftpPool.Delete(targetName)
	}

	cli, err := sftp.NewClient(targetClient)
	if err != nil {
		return nil, err
	}

	sftpPool.Store(targetName, cli)
	return cli, nil
}

func startSFTPServer(
	channel ssh.Channel,
	targetClient *ssh.Client,
	targetName string,
) {
	defer channel.Close()

	targetSftp, err := getSftpClient(targetName, targetClient)
	if err != nil {
		global.GVA_LOG.Error("get target sftp client failed: ", zap.Any("jump server", err))
		return
	}

	handler := &ProxySFTPHandler{
		client: targetSftp,
	}

	server := sftp.NewRequestServer(
		channel,
		sftp.Handlers{
			FileGet:  handler,
			FilePut:  handler,
			FileCmd:  handler,
			FileList: handler,
		},
	)

	global.GVA_LOG.Info("sftp proxy started for ", zap.String("target", targetName))

	if err = server.Serve(); err != nil && err != io.EOF {
		global.GVA_LOG.Error("sftp serve error: ", zap.Any("jump server", err))
	}
}
