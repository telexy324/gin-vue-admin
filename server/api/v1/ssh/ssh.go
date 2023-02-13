package ssh

import (
	"fmt"
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

func ShellWeb(c echo.Context) error {
	var err error

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		ubzer.MLog.Error("websocket upgrade 失败", zap.Error(err))
		return err
	}
	_, readContent, err := conn.ReadMessage()
	if err != nil {
		ubzer.MLog.Error("websocket 读取ip、用户名、密码 失败", zap.Error(err))
		return err
	}
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ readContent: %v\n", string(readContent))

	sshClient, err := connections.DecodeMsgToSSHClient(string(readContent))
	if err != nil {
		return err
	}
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ sshClient: %v\n", sshClient)

	terminal := connections.Terminal{
		Columns: 150,
		Rows:    35,
	}

	var port = 22
	err = sshClient.GenerateClient(sshClient.IpAddress, sshClient.Username, sshClient.Password, port)
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		return err
	}
	sshClient.RequestTerminal(terminal)
	sshClient.Connect(conn)
	return nil
}
