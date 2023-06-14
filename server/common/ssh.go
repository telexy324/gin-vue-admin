package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

const _hex = "0123456789abcdef"

type SshService struct {
}

type PtyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	ModeList string
}

type Terminal struct {
	Columns uint32 `json:"cols"`
	Rows    uint32 `json:"rows"`
}

type SSHClient struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ManageIp string `json:"manageIp"`
	SshPort  int    `json:"sshPort"`
	//Session  *ssh.Session
	Client *ssh.Client
	//Channel ssh.Channel
}

func newSSHClient() SSHClient {
	client := SSHClient{}
	return client
}

const (
	wsMsgCmd    = "cmd"
	wsMsgResize = "resize"
)

type wsMsg struct {
	Type string `json:"type"`
	Cmd  string `json:"cmd"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

// connect to ssh server using ssh session.
type SshConn struct {
	// calling Write() to write data into ssh server
	StdinPipe io.WriteCloser
	// Write() be called to receive data from ssh server
	ComboOutput *wsBufferWriter
	Session     *ssh.Session
}

type wsBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

// implement Write interface to write bytes from ssh server into bytes.Buffer.
func (w *wsBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

// flushComboOutput flush ssh.session combine output into websocket response
func flushComboOutput(w *wsBufferWriter, wsConn *websocket.Conn) error {
	if w.buffer.Len() != 0 {
		err := wsConn.WriteMessage(websocket.TextMessage, w.buffer.Bytes())
		if err != nil {
			return err
		}
		//global.GVA_LOG.Info("ssh direct message ", zap.ByteString("", w.buffer.Bytes()))
		w.buffer.Reset()
	}
	return nil
}

func DecodeMsgToSSHClient(msg string) (SSHClient, error) {
	client := newSSHClient()
	decoded, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return client, err
	}
	err = json.Unmarshal(decoded, &client)
	if err != nil {
		return client, err
	}
	return client, nil
}

func FillSSHClient(ip, username, password string, port int) (SSHClient, error) {
	client := newSSHClient()
	client.ManageIp, client.SshPort = ip, port
	client.Username, client.Password = username, password
	return client, nil
}

func (c *SSHClient) GenerateClient() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	if len(c.Password) > 0 {
		auth = append(auth, ssh.Password(c.Password))
	} else {
		homePath, e := os.UserHomeDir()
		if e != nil {
			return e
		}
		key, e := os.ReadFile(path.Join(homePath, ".ssh", "id_rsa"))
		if e != nil {
			return e
		}
		signer, e := ssh.ParsePrivateKey(key)
		if e != nil {
			return e
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:    c.Username,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		//HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		//	return nil
		//},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr = fmt.Sprintf("%s:%d", c.ManageIp, c.SshPort)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	c.Client = client
	return nil
}

func (c *SSHClient) Command(command string, logger Logger, manageIP string) (err error) {
	logger.Log(command, manageIP)
	session, err := c.Client.NewSession()
	if err != nil {
		logger.Log("ssh open session failed")
		return
	}
	output, err := session.CombinedOutput(command)
	if err != nil {
		logger.Log("ssh exec failed")
		return
	}
	logger.Log(string(output), manageIP)
	return
}

func (c *SSHClient) CommandSingle(command string) (output string, err error) {
	session, err := c.Client.NewSession()
	if err != nil {
		return
	}
	outputBytes, err := session.CombinedOutput(command)
	if err != nil {
		return
	}
	return string(outputBytes), nil
}

type singleWriter struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (w *singleWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (c *SSHClient) Commands(command string, logger Logger, manageIP string) (err error) {
	logger.Log(command, manageIP)
	session, err := c.Client.NewSession()
	if err != nil {
		logger.Log("ssh open session failed")
		return
	}
	//stdinP, err := session.StdinPipe()
	//if err != nil {
	//	logger.Log("ssh open session stdin failed")
	//	return
	//}
	//buffer := new(bytes.Buffer)
	var buffer singleWriter
	session.Stdout = &buffer
	session.Stderr = &buffer
	//if err = session.Shell(); err != nil {
	//	global.GVA_LOG.Info("ssh session shell error ", zap.Any("err ", err))
	//	return
	//}
	quitChan := make(chan bool)
	go func() {
		//every 120ms write combine output bytes into websocket response
		tick := time.NewTicker(time.Millisecond * time.Duration(100))
		//for range time.Tick(120 * time.Millisecond){}
		defer tick.Stop()
		var last = 0
		for {
			select {
			case <-tick.C:
				//if outputBuffer.Len() != 0 {
				//	logger.Log(outputBuffer.String(), manageIP)
				//	outputBuffer.Reset()
				//}
				//if outputBuffer.Len() != 0 {
				//	logger.Log(string(outputBuffer.Next(last)), manageIP)
				//	last = outputBuffer.Len()
				//	outputBuffer.Reset()
				//}
				//buf:=make([]byte,0,4096)
				//outputBuffer.Write(buf)
				//logger.Log(string(buf), manageIP)
				if buffer.b.Len() != 0 {
					last = buffer.b.Len()
					logger.Log(string(buffer.b.Next(last)), manageIP)
				}
			case <-quitChan:
				//global.GVA_LOG.Info("ssh receive quit")
				//logger.Log(string(buffer.b.Next(last)), manageIP)
				return
			}
		}
	}()
	//go func(sess *ssh.Session, exitCh chan bool) {
	//	if err = sess.Wait(); err != nil {
	//		global.GVA_LOG.Info("ssh receive error ", zap.Any("err ", err))
	//		exitCh <- true
	//	}
	//}(session, quitChan)

	//if err = session.Start(command); err != nil {
	//	global.GVA_LOG.Info("command ", zap.Any("command ", command))
	//	global.GVA_LOG.Error("ws cmd bytes write to ssh.stdin pipe failed", zap.Any("err ", err))
	//}
	//if err = session.Wait(); err != nil {
	//	global.GVA_LOG.Info("ssh receive error ", zap.Any("err ", err))
	//	quitChan <- true
	//}
	err = session.Run(command)
	if err != nil {
		global.GVA_LOG.Error("ws cmd bytes write to ssh.stdin pipe failed", zap.Any("err ", err))
	}
	time.Sleep(time.Millisecond * time.Duration(250))
	quitChan <- true
	return
}

func (c *SSHClient) NewSftp(opts ...sftp.ClientOption) (*sftp.Client, error) {
	return sftp.NewClient(c.Client, opts...)
}

func (c *SSHClient) Upload(file io.Reader, remotePath string) (err error) {
	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Create(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	_, err = io.Copy(remote, file)
	return
}

// Download file from remote server!
func (c *SSHClient) Download(remotePath string) (fileBytes []byte, err error) {
	//
	//local, err := os.Create(localPath)
	//if err != nil {
	//	return
	//}
	//defer local.Close()

	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Open(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	return io.ReadAll(remote)
	//if _, err = io.Copy(local, remote); err != nil {
	//	return
	//}
	//
	//return local.Sync()
	//return remote, nil
}

// Download file from remote server!
func (c *SSHClient) DownloadReader(remotePath string) (file io.Reader, err error) {
	ftp, err := c.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Open(remotePath)
	if err != nil {
		return
	}
	defer remote.Close()

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, remote); err != nil {
		return
	}
	return buf, nil
}

//func (c *SSHClient) RequestTerminal(terminal Terminal) *SSHClient {
//	//session, err := c.Client.NewSession()
//	//if err != nil {
//	//	return nil
//	//}
//	//c.Session = session
//	channel, inRequests, err := c.Client.OpenChannel("session", nil)
//	if err != nil {
//		return nil
//	}
//	c.Channel = channel
//	go ssh.DiscardRequests(inRequests)
//	//go func() {
//	//	for req := range inRequests {
//	//		if req.WantReply {
//	//			req.Reply(false, nil)
//	//		}
//	//	}
//	//}()
//	modes := ssh.TerminalModes{
//		ssh.ECHO:          1,
//		ssh.TTY_OP_ISPEED: 14400,
//		ssh.TTY_OP_OSPEED: 14400,
//	}
//	var modeList []byte
//	for k, v := range modes {
//		kv := struct {
//			Key byte
//			Val uint32
//		}{k, v}
//		modeList = append(modeList, ssh.Marshal(&kv)...)
//	}
//	modeList = append(modeList, 0)
//	req := PtyRequestMsg{
//		Term:     "xterm",
//		Columns:  terminal.Columns,
//		Rows:     terminal.Rows,
//		Width:    terminal.Columns * 8,
//		Height:   terminal.Columns * 8,
//		ModeList: string(modeList),
//	}
//	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
//	if !ok || err != nil {
//		return nil
//	}
//	ok, err = channel.SendRequest("shell", true, nil)
//	if !ok || err != nil {
//		return nil
//	}
//	return c
//}
//
//func (c *SSHClient) Connect(ws *websocket.Conn) {
//	go func() {
//		for {
//			_, p, err := ws.ReadMessage()
//			if err != nil {
//				return
//			}
//			_, err = c.Channel.Write(p)
//			if err != nil {
//				return
//			}
//		}
//	}()
//
//	go func() {
//		br := bufio.NewReader(c.Channel)
//		buf := []byte{}
//		t := time.NewTimer(time.Microsecond * 100)
//		defer t.Stop()
//		r := make(chan rune)
//
//		go func() {
//			defer c.Client.Close()
//			// defer c.Client.Close()
//
//			for {
//				x, size, err := br.ReadRune()
//				if err != nil {
//					log.Println(err)
//					ws.WriteMessage(1, []byte("\033[31m已经关闭连接!\033[0m"))
//					ws.Close()
//					return
//				}
//				if size > 0 {
//					r <- x
//				}
//			}
//		}()
//
//		for {
//			select {
//			case <-t.C:
//				if len(buf) != 0 {
//					err := ws.WriteMessage(websocket.TextMessage, buf)
//					buf = []byte{}
//					if err != nil {
//						log.Println(err)
//						return
//					}
//				}
//				t.Reset(time.Microsecond * 100)
//			case d := <-r:
//				if d != utf8.RuneError {
//					p := make([]byte, utf8.RuneLen(d))
//					utf8.EncodeRune(p, d)
//					buf = append(buf, p...)
//				} else {
//					buf = append(buf, []byte("@")...)
//				}
//			}
//		}
//	}()
//
//	defer func() {
//		if err := recover(); err != nil {
//			log.Println(err)
//		}
//	}()
//}

//func (c *SSHClient) RequestShell() *SSHClient {
//	channel, inRequests, err := c.Client.OpenChannel("session", nil)
//	if err != nil {
//		return nil
//	}
//	c.Channel = channel
//	go ssh.DiscardRequests(inRequests)
//	ok, err := channel.SendRequest("shell", true, nil)
//	if !ok || err != nil {
//		return nil
//	}
//	return c
//}
//
//func (c *SSHClient) ConnectShell(shell string, logger common.Logger) (err error) {
//	logger.LogSsh(c.Channel)
//	_, err = c.Channel.Write([]byte(shell))
//	if err != nil {
//		return
//	}
//
//	//defer func() {
//	//	if err := recover(); err != nil {
//	//		log.Println(err)
//	//	}
//	//}()
//	return
//}

// setup ssh shell session
// set Session and StdinPipe here,
// and the Session.Stdout and Session.Sdterr are also set.
func NewSshConn(sshClient *ssh.Client, cols, rows int) (*SshConn, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	// we set stdin, then we can write data to ssh server via this stdin.
	// but, as for reading data from ssh server, we can set Session.Stdout and Session.Stderr
	// to receive data from ssh server, and write back to somewhere.
	stdinP, err := sshSession.StdinPipe()
	if err != nil {
		return nil, err
	}

	comboWriter := new(wsBufferWriter)
	//ssh.stdout and stderr will write output into comboWriter
	sshSession.Stdout = comboWriter
	sshSession.Stderr = comboWriter

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err = sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil, err
	}
	// Start remote shell
	if err = sshSession.Shell(); err != nil {
		return nil, err
	}
	return &SshConn{StdinPipe: stdinP, ComboOutput: comboWriter, Session: sshSession}, nil
}

func (s *SshConn) Close() {
	if s.Session != nil {
		s.Session.Close()
	}

}

// ReceiveWsMsg  receive websocket msg do some handling then write into ssh.session.stdin
func (ssConn *SshConn) ReceiveWsMsg(wsConn *websocket.Conn, logBuff *bytes.Buffer, exitCh chan bool) {
	//tells other go routine quit
	defer setQuit(exitCh)
	for {
		select {
		case <-exitCh:
			return
		default:
			//read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				global.GVA_LOG.Error("reading webSocket message failed", zap.Any("err ", err))
				return
			}
			//unmashal bytes into struct
			msgObj := wsMsg{
				Type: "cmd",
				Cmd:  "",
				Rows: 50,
				Cols: 180,
			}
			//if err := json.Unmarshal(wsData, &msgObj); err != nil {
			//	logrus.WithError(err).WithField("wsData", string(wsData)).Error("unmarshal websocket message failed")
			//}
			switch msgObj.Type {
			case wsMsgResize:
				//handle xterm.js size change
				if msgObj.Cols > 0 && msgObj.Rows > 0 {
					if err = ssConn.Session.WindowChange(msgObj.Rows, msgObj.Cols); err != nil {
						global.GVA_LOG.Error("ssh pty change windows size failed", zap.Any("err ", err))
					}
				}
			case wsMsgCmd:
				//handle xterm.js stdin
				//decodeBytes, err := base64.StdEncoding.DecodeString(msgObj.Cmd)
				decodeBytes := wsData
				if err != nil {
					global.GVA_LOG.Error("websock cmd string base64 decoding failed", zap.Any("err ", err))
				}
				if _, err = ssConn.StdinPipe.Write(decodeBytes); err != nil {
					global.GVA_LOG.Error("ws cmd bytes write to ssh.stdin pipe failed", zap.Any("err ", err))
				}
				//write input cmd to log buffer
				if _, err = logBuff.Write(decodeBytes); err != nil {
					global.GVA_LOG.Error("write received cmd into log buffer failed", zap.Any("err ", err))
				}
			}
		}
	}
}

func (ssConn *SshConn) SendComboOutput(wsConn *websocket.Conn, exitCh chan bool) {
	//tells other go routine quit
	defer setQuit(exitCh)

	//every 120ms write combine output bytes into websocket response
	tick := time.NewTicker(time.Millisecond * time.Duration(120))
	//for range time.Tick(120 * time.Millisecond){}
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			//write combine output bytes into websocket response
			if err := flushComboOutput(ssConn.ComboOutput, wsConn); err != nil {
				global.GVA_LOG.Error("ssh sending combo output to webSocket failed", zap.Any("err ", err))
				return
			}
		case <-exitCh:
			return
		}
	}
}

func (ssConn *SshConn) SessionWait(quitChan chan bool) {
	if err := ssConn.Session.Wait(); err != nil {
		global.GVA_LOG.Error("ssh session wait failed", zap.Any("err ", err))
		setQuit(quitChan)
	}
}

func setQuit(ch chan bool) {
	ch <- true
}

//	func (ssConn *SshConn) SendCommand(command string, logger common.Logger, exitCh chan bool) (err error) {
//		logger.Log(command)
//		_, err = ssConn.StdinPipe.Write([]byte(command))
//		if err != nil {
//			logger.Log("ssh run command failed")
//			return
//		}
//		tick := time.NewTicker(time.Millisecond * time.Duration(120))
//		//for range time.Tick(120 * time.Millisecond){}
//		defer tick.Stop()
//		for {
//			select {
//			case <-tick.C:
//				logger.Log(string(ssConn.ComboOutput.buffer.Bytes()))
//				ssConn.ComboOutput.buffer.Reset()
//			case <-exitCh:
//				return
//			}
//		}
//	}

func (c *SSHClient) CommandBatch(commands []string, logger Logger, manageIP string) (err error) {
	session, err := c.Client.NewSession()
	if err != nil {
		global.GVA_LOG.Error("ssh open session failed ", zap.Any("err ", err))
		logger.Log("ssh open session failed")
		return
	}
	defer session.Close()

	stdinP, err := session.StdinPipe()
	if err != nil {
		global.GVA_LOG.Error("ssh open stdin failed ", zap.Any("err ", err))
		logger.Log("ssh open stdin failed")
		return
	}

	buffer := new(singleWriter)
	session.Stdout = buffer
	session.Stderr = buffer

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err = session.RequestPty("xterm", 150, 35, modes); err != nil {
		global.GVA_LOG.Error("ssh request pty failed ", zap.Any("err ", err))
		logger.Log("ssh request pty failed")
		return
	}
	// Start remote shell
	if err = session.Shell(); err != nil {
		global.GVA_LOG.Error("ssh start shell failed ", zap.Any("err ", err))
		logger.Log("ssh start shell failed")
		return
	}

	quitChan := make(chan bool)
	go func() {
		//every 120ms write combine output bytes into websocket response
		tick := time.NewTicker(time.Millisecond * time.Duration(100))
		//for range time.Tick(120 * time.Millisecond){}
		defer tick.Stop()
		//var last = 0
		for {
			select {
			case <-tick.C:
				//if buffer.b.Len() != 0 {
				//	last = buffer.b.Len()
				//	logger.Log(string(buffer.b.Next(last)), manageIP)
				//}
				//if buffer.b.Len() != 0 {
				//	raw := buffer.b.Bytes()
				//	global.GVA_LOG.Info(string(raw), zap.ByteString("", raw))
				//	processed := safeAddByteString(raw)
				//	logger.Log(string(processed), manageIP)
				//}
				//buffer.b.Reset()
				if buffer.b.Len() != 0 {
					rawString := string(buffer.b.Bytes())
					cleanMsg := stripansi.Strip(rawString)
					logger.Log(cleanMsg, manageIP)
				}
				buffer.b.Reset()
				//if buffer.b.Len() != 0 {
				//	//rawString := string(buffer.b.Bytes())
				//	//raws := strings.Split(rawString, "\r\n")
				//	//for _, raw := range raws {
				//	//	if !(strings.Contains(raw, "Last login") || strings.Contains(raw, ` ~]# `)) {
				//	//		logger.Log(string(buffer.b.Bytes()), manageIP)
				//	//	}
				//	//}
				//	raw := make([]rune, 0)
				//	var through bool
				//	for {
				//		r, _, err := buffer.b.ReadRune()
				//		if err != nil && err != io.EOF {
				//			global.GVA_LOG.Error("ssh read rune failed ", zap.Any("err ", err))
				//			through = true
				//			break
				//		} else if err != nil && err == io.EOF {
				//			break
				//		}
				//		//else if r == []rune("\u0007")[0] {
				//		//	through = true
				//		//}
				//		raw = append(raw, r)
				//	}
				//	if through {
				//		continue
				//	}
				//	rawStrings := strings.Split(string(raw), "\r\n")
				//	for _, rs := range rawStrings {
				//		if !(strings.Contains(rs, "Last login")) {
				//			logger.Log(rs, manageIP)
				//		}
				//	}
				//	buffer.b.Reset()
				//}
			case <-quitChan:
				return
			}
		}
	}()
	go func() {
		for _, command := range commands {
			logger.Log(command, manageIP)
			if _, err = stdinP.Write([]byte(command + "\n")); err != nil {
				global.GVA_LOG.Error("cmd bytes write to ssh.stdin pipe failed", zap.Any("err ", err))
			}
			time.Sleep(time.Millisecond * time.Duration(200))
		}
		_, err = stdinP.Write([]byte("exit" + "\n"))
		if err != nil {
			global.GVA_LOG.Error("ssh exit failed", zap.Any("err ", err))
		}
	}()
	err = session.Wait()
	if err != nil {
		global.GVA_LOG.Info("ssh receive error ", zap.Any("err ", err))
	}
	time.Sleep(time.Millisecond * time.Duration(300))
	quitChan <- true
	return
}

//func (c *SSHClient) CommandSerial(commands []string, logger Logger, manageIP string) (err error) {
//	session, err := c.Client.NewSession()
//	if err != nil {
//		global.GVA_LOG.Error("ssh open session failed ", zap.Any("err ", err))
//		logger.Log("ssh open session failed")
//		return
//	}
//	defer session.Close()
//
//	modes := ssh.TerminalModes{
//		ssh.ECHO:          1,     // disable echo
//		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
//		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
//	}
//	// Request pseudo terminal
//	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
//		global.GVA_LOG.Error("ssh request pty failed ", zap.Any("err ", err))
//		logger.Log("ssh request pty failed")
//		return
//	}
//	w, err := session.StdinPipe()
//	if err != nil {
//		global.GVA_LOG.Error("ssh get stdin failed ", zap.Any("err ", err))
//		logger.Log("ssh get stdin failed")
//		return
//	}
//	r, err := session.StdoutPipe()
//	if err != nil {
//		global.GVA_LOG.Error("ssh get stdout failed ", zap.Any("err ", err))
//		logger.Log("ssh get stdout failed")
//		return
//	}
//	if err = session.Start("/bin/sh"); err != nil {
//		global.GVA_LOG.Error("ssh start failed ", zap.Any("err ", err))
//		logger.Log("ssh start failed")
//		return
//	}
//	readUntil(r, escapePrompt) //ignore the shell output
//	for _, command := range commands {
//		write(w, command)
//		out, _ := readUntil(r, escapePrompt)
//		logger.Log(*out, manageIP)
//	}
//	write(w, "exit")
//
//	session.Wait()
//	return
//}
//
//var escapePrompt = []byte{'$', ' '}
//
//func write(w io.WriteCloser, command string) error {
//	_, err := w.Write([]byte(command + "\n"))
//	return err
//}
//
//func readUntil(r io.Reader, matchingByte []byte) (*string, error) {
//	var buf [64 * 1024]byte
//	var t int
//	for {
//		n, err := r.Read(buf[t:])
//		if err != nil {
//			return nil, err
//		}
//		t += n
//		if isMatch(buf[:t], t, matchingByte) {
//			stringResult := string(buf[:t])
//			return &stringResult, nil
//		}
//	}
//}
//
//func isMatch(bytes []byte, t int, matchingBytes []byte) bool {
//	if t >= len(matchingBytes) {
//		for i := 0; i < len(matchingBytes); i++ {
//			if bytes[t-len(matchingBytes)+i] != matchingBytes[i] {
//				return false
//			}
//		}
//		return true
//	}
//	return false
//}

//func safeAddByteString(s []byte) (out []byte) {
//	out = make([]byte, 0, 8)
//	for i := 0; i < len(s); {
//		if tryAddRuneSelf(s[i], &out) {
//			i++
//			continue
//		}
//		r, size := utf8.DecodeRune(s[i:])
//		if tryAddRuneError(r, size, &out) {
//			i++
//			continue
//		}
//		out = append(out, s[i:i+size]...)
//		i += size
//	}
//	return out
//}
//
//func tryAddRuneSelf(b byte, out *[]byte) bool {
//	if b >= utf8.RuneSelf {
//		return false
//	}
//	if 0x20 <= b && b != '\\' && b != '"' {
//		appendByte(b, out)
//		return true
//	}
//	switch b {
//	case '\\', '"':
//		appendByte('\\', out)
//		appendByte(b, out)
//	case '\n':
//		appendByte('\\', out)
//		appendByte('n', out)
//	case '\r':
//		appendByte('\\', out)
//		appendByte('r', out)
//	case '\t':
//		appendByte('\\', out)
//		appendByte('t', out)
//	default:
//		// Encode bytes < 0x20, except for the escape sequences above.
//		appendString(`\u00`, out)
//		appendByte(_hex[b>>4], out)
//		appendByte(_hex[b&0xF], out)
//	}
//	return true
//}
//
//func tryAddRuneError(r rune, size int, out *[]byte) bool {
//	if r == utf8.RuneError && size == 1 {
//		appendString(`\ufffd`, out)
//		return true
//	}
//	return false
//}
//
//// AppendByte writes a single byte to the Buffer.
//func appendByte(v byte, out *[]byte) {
//	*out = append(*out, v)
//}
//
//// AppendString writes a string to the Buffer.
//func appendString(s string, out *[]byte) {
//	*out = append(*out, s...)
//}
