package ssh

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	ssh2 "github.com/flipped-aurora/gin-vue-admin/server/service/ssh"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

type SshApi struct {
}

var (
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (a *SshApi) ShellWeb(c *gin.Context) error {
	var err error

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("websocket upgrade 失败", zap.Error(err))
		return err
	}
	_, readContent, err := conn.ReadMessage()
	if err != nil {
		global.GVA_LOG.Error("websocket 读取ip、用户名、密码 失败", zap.Error(err))
		return err
	}
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ readContent: %v\n", string(readContent))

	sshClient, err := sshService.DecodeMsgToSSHClient(string(readContent))
	if err != nil {
		return err
	}
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ sshClient: %v\n", sshClient)

	terminal := ssh2.Terminal{
		Columns: 150,
		Rows:    35,
	}

	err, server := cmdbServerService.GetServerById(float64(sshClient.Server.ID))
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		return err
	}
	sshClient.Server = &server

	err = sshClient.GenerateClient(sshClient.Server.ManageIp, sshClient.Username, sshClient.Password, sshClient.Server.SshPort)
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		return err
	}
	sshClient.RequestTerminal(terminal)
	sshClient.Connect(conn)
	return nil
}
