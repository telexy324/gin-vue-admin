package jumpServer

import (
	"context"
	"fmt"
	"io"
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
	Parent       *Session // null for parent
	StartedAt    time.Time
	LastActiveAt atomic.Int64

	Ctx    context.Context
	Cancel context.CancelFunc

	children sync.Map // childID -> *Session
}

type activityWriter struct {
	rw      io.ReadWriter
	session *Session
}

func (a *activityWriter) Write(p []byte) (int, error) {
	n, err := a.rw.Write(p)
	a.session.LastActiveAt.Store(time.Now().Unix())
	if err != nil {
		a.session.End()
	}
	return n, err
}

func (a *activityWriter) Read(p []byte) (int, error) {
	n, err := a.rw.Read(p)
	a.session.LastActiveAt.Store(time.Now().Unix())
	if err != nil {
		a.session.End()
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
	idleTimeout    = 10 * time.Minute
	maxSessionTime = 2 * time.Hour
)

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
	parent := newParentSession("", targetName)

	// 3. 连接目标服务器
	targetClient, err := connectTarget(targetName)
	if err != nil {
		global.GVA_LOG.Error("target connect failed: ", zap.Any("jump server", err))
		return
	}
	defer targetClient.Close()

	go ssh.DiscardRequests(reqs)

	for newCh := range chans {
		child := parent.NewChild()
		go handleNewChannel(newCh, targetClient, child)
	}

}

func connectTarget(targetName string) (*ssh.Client, error) {
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

	return ssh.Dial("tcp", addr, cfg)
}

func newParentSession(user, target string) *Session {
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

	for {
		select {
		case <-sess.Ctx.Done():
			return

		case <-ticker.C:
			now := time.Now()

			// 最大会话时长
			if now.Sub(sess.StartedAt) > maxSessionTime {
				global.GVA_LOG.Info("session max time reached: ", zap.String("jump server", sess.ID))
				sess.End()
				return
			}

			// 空闲超时
			last := time.Unix(sess.LastActiveAt.Load(), 0)
			if now.Sub(last) > idleTimeout {
				global.GVA_LOG.Info("session idle timeout: ", zap.String("jump server", sess.ID))
				sess.End()
				return
			}
		}
	}
}

func getSftpClient(targetName string, targetClient *ssh.Client) (*sftp.Client, error) {
	cli, err := sftp.NewClient(targetClient)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func startSFTPServer(
	channel ssh.Channel,
	targetClient *ssh.Client,
	sess *Session,
) {
	defer channel.Close()

	targetSftp, err := getSftpClient(sess.Target, targetClient)
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

	global.GVA_LOG.Info("sftp proxy started for ", zap.String("target", sess.Target))

	if err := server.Serve(); err != nil && err != io.EOF {
		global.GVA_LOG.Error("sftp serve error: ", zap.Any("jump server", err))
	}
	sess.End()
}

func handleNewChannel(
	newChannel ssh.NewChannel,
	targetClient *ssh.Client,
	sess *Session,
) {
	if newChannel.ChannelType() != "session" {
		newChannel.Reject(ssh.UnknownChannelType, "unsupported channel")
		return
	}

	// 接受源 channel
	srcChannel, srcRequests, err := newChannel.Accept()
	if err != nil {
		global.GVA_LOG.Error("failed to accept channel: ", zap.Any("jump server", err))
		return
	}

	// 每个源 channel 对应独立 dstSession （关键！支持 clone）
	dstSession, err := targetClient.NewSession()
	if err != nil {
		global.GVA_LOG.Error("create session failed: ", zap.Any("jump server", err))
		srcChannel.Close()
		return
	}

	// 会话关闭时的清理
	go func() {
		<-sess.Ctx.Done()
		srcChannel.Close()
		dstSession.Close()
	}()

	go func() {
		for req := range srcRequests {

			switch req.Type {

			case "pty-req":
				dstSession.RequestPty("xterm-256color", 40, 120, ssh.TerminalModes{})
				req.Reply(true, nil)

			case "shell":
				dstSession.Stdout = &activityWriter{rw: srcChannel, session: sess}
				dstSession.Stderr = &activityWriter{rw: srcChannel, session: sess}
				dstSession.Stdin = &activityWriter{rw: srcChannel, session: sess}

				dstSession.Shell()

				go func() {
					err := dstSession.Wait()
					global.GVA_LOG.Error("target session exited: ", zap.Any("jump server", err))
					sess.End()
				}()

				req.Reply(true, nil)

			case "subsystem":
				if string(req.Payload[4:]) == "sftp" {
					req.Reply(true, nil)
					startSFTPServer(srcChannel, targetClient, sess)
					return
				}
				req.Reply(false, nil)
			}
		}
	}()
	go monitorSession(sess)
}

func (s *Session) NewChild() *Session {
	ctx, cancel := context.WithCancel(s.Ctx)

	child := &Session{
		ID:        uuid.NewString(),
		User:      s.User,
		Target:    s.Target,
		StartedAt: time.Now(),
		Ctx:       ctx,
		Cancel:    cancel,
		Parent:    s,
	}

	child.LastActiveAt.Store(time.Now().Unix())
	s.children.Store(child.ID, child)
	return child
}

func (s *Session) RemoveChild(id string) {
	s.children.Delete(id)
	// 只有父会话在没有 child 时才退出
	if s.Parent == nil {
		empty := true
		s.children.Range(func(_, _ any) bool {
			empty = false
			return false
		})
		if empty {
			s.Cancel()
		}
	}
}

func (s *Session) End() {
	if s.Parent != nil {
		s.Parent.RemoveChild(s.ID)
	}
	s.Cancel()
}
