package ssh

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/taskPool"
	"log"
	"net"
	"os"
	"path"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
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
	Client  *ssh.Client
	channel ssh.Channel
}

func newSSHClient() SSHClient {
	client := SSHClient{}
	return client
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
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", c.Server.ManageIp, c.Server.SshPort)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	c.Client = client
	return nil
}

func (c *SSHClient) RequestTerminal(terminal Terminal) *SSHClient {
	//session, err := c.Client.NewSession()
	//if err != nil {
	//	return nil
	//}
	//c.Session = session
	channel, inRequests, err := c.Client.OpenChannel("session", nil)
	if err != nil {
		return nil
	}
	c.channel = channel
	go ssh.DiscardRequests(inRequests)
	//go func() {
	//	for req := range inRequests {
	//		if req.WantReply {
	//			req.Reply(false, nil)
	//		}
	//	}
	//}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := PtyRequestMsg{
		Term:     "xterm",
		Columns:  terminal.Columns,
		Rows:     terminal.Rows,
		Width:    terminal.Columns * 8,
		Height:   terminal.Columns * 8,
		ModeList: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		return nil
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		return nil
	}
	return c
}

func (c *SSHClient) Connect(ws *websocket.Conn) {
	go func() {
		for {
			_, p, err := ws.ReadMessage()
			if err != nil {
				return
			}
			_, err = c.channel.Write(p)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		br := bufio.NewReader(c.channel)
		buf := []byte{}
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		r := make(chan rune)

		go func() {
			defer c.Client.Close()
			// defer c.Client.Close()

			for {
				x, size, err := br.ReadRune()
				if err != nil {
					log.Println(err)
					ws.WriteMessage(1, []byte("\033[31m已经关闭连接!\033[0m"))
					ws.Close()
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		for {
			select {
			case <-t.C:
				if len(buf) != 0 {
					err := ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						log.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
}

func (c *SSHClient) RequestShell() *SSHClient {
	channel, inRequests, err := c.Client.OpenChannel("session", nil)
	if err != nil {
		return nil
	}
	c.channel = channel
	go ssh.DiscardRequests(inRequests)
	ok, err := channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		return nil
	}
	return c
}

func (c *SSHClient) ConnectShell(shell string, taskRunner taskPool.TaskRunner) (err error) {
	_, err = c.channel.Write([]byte(shell))
	if err != nil {
		return
	}

	taskRunner.LogSsh(&c.channel)

	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Println(err)
	//	}
	//}()
	return
}
