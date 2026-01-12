package jumpServer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/sftp"
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

//func (a *activityWriter) Write(p []byte) (int, error) {
//	a.session.LastActiveAt.Store(time.Now().Unix())
//	return a.rw.Write(p)
//}
//
//func (a *activityWriter) Read(p []byte) (int, error) {
//	a.session.LastActiveAt.Store(time.Now().Unix())
//	return a.rw.Read(p)
//}

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
	idleTimeout    = 10 * time.Minute
	maxSessionTime = 2 * time.Hour
)

func main() {
	// 1. SSH Server 配置
	privateBytes, err := os.ReadFile("server_host_key")
	if err != nil {
		log.Fatal(err)
	}
	private, _ := ssh.ParsePrivateKey(privateBytes)

	config := &ssh.ServerConfig{
		NoClientAuth: true, // 简化示例（生产请做认证）
	}
	config.AddHostKey(private)

	// 2. 监听 22（或其他端口）
	listener, err := net.Listen("tcp", ":44488")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Jump server listening on :44488")

	for {
		conn, _ := listener.Accept()
		go handleConn(conn, config)
	}
}

func handleConn(nConn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Println("handshake failed:", err)
		return
	}
	defer sshConn.Close()

	targetName := sshConn.User()
	sess := newSession("", targetName)

	// 3. 连接目标服务器
	targetClient, err := connectTarget(targetName)
	if err != nil {
		log.Println("target connect failed:", err)
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
						log.Println("starting sftp subsystem")
						req.Reply(true, nil)

						startSFTPServer(srcChannel, targetClient)
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
	}
	addr := fmt.Sprintf("%s:%d", targetName, 1122)
	return ssh.Dial("tcp", addr, cfg)
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

	for {
		select {
		case <-sess.Ctx.Done():
			return

		case <-ticker.C:
			now := time.Now()

			// 最大会话时长
			if now.Sub(sess.StartedAt) > maxSessionTime {
				log.Println("session max time reached:", sess.ID)
				sess.Cancel()
				return
			}

			// 空闲超时
			last := time.Unix(sess.LastActiveAt.Load(), 0)
			if now.Sub(last) > idleTimeout {
				log.Println("session idle timeout:", sess.ID)
				sess.Cancel()
				return
			}
		}
	}
}

func startSFTPServer(
	channel ssh.Channel,
	targetClient *ssh.Client,
) {
	defer channel.Close()
	defer targetClient.Close()

	targetSftp, err := sftp.NewClient(targetClient)
	if err != nil {
		log.Println("create target sftp client failed:", err)
		return
	}
	defer targetSftp.Close()

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

	log.Println("sftp proxy started")

	if err := server.Serve(); err != nil && err != io.EOF {
		log.Println("sftp serve error:", err)
	}

	log.Println("sftp proxy closed")
}
