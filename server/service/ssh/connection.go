package ssh

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

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
	Username string                         `json:"username"`
	Password string                         `json:"password"`
	Server   *application.ApplicationServer `json:"server"`
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
		w.buffer.Reset()
	}
	return nil
}

func (sshService *SshService) DecodeMsgToSSHClient(msg string) (SSHClient, error) {
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

func (sshService *SshService) FillSSHClient(ip, username, password string, port int) (SSHClient, error) {
	client := newSSHClient()
	client.Server = &application.ApplicationServer{
		ManageIp: ip,
		SshPort:  port,
	}
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
	addr = fmt.Sprintf("%s:%d", c.Server.ManageIp, c.Server.SshPort)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	c.Client = client
	return nil
}

func (c *SSHClient) Command(command string, logger common.Logger) (err error) {
	logger.Log(command)
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
	logger.Log(string(output))
	return
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
func (sshService *SshService) NewSshConn(sshClient *ssh.Client, cols, rows int) (*SshConn, error) {
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

//func (ssConn *SshConn) SendCommand(command string, logger common.Logger, exitCh chan bool) (err error) {
//	logger.Log(command)
//	_, err = ssConn.StdinPipe.Write([]byte(command))
//	if err != nil {
//		logger.Log("ssh run command failed")
//		return
//	}
//	tick := time.NewTicker(time.Millisecond * time.Duration(120))
//	//for range time.Tick(120 * time.Millisecond){}
//	defer tick.Stop()
//	for {
//		select {
//		case <-tick.C:
//			logger.Log(string(ssConn.ComboOutput.buffer.Bytes()))
//			ssConn.ComboOutput.buffer.Reset()
//		case <-exitCh:
//			return
//		}
//	}
//}
