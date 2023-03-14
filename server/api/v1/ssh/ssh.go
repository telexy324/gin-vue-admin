package ssh

import (
	"bytes"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
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

func (a *SshApi) ShellWeb(c *gin.Context) {
	var err error

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("websocket upgrade 失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	_, readContent, err := conn.ReadMessage()
	fmt.Println(err)
	if err != nil {
		global.GVA_LOG.Error("websocket 读取ip、用户名、密码 失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ readContent: %v\n", string(readContent))

	sshClient, err := common.DecodeMsgToSSHClient(string(readContent))
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ sshClient: %v\n", sshClient)

	//terminal := ssh2.Terminal{
	//	Columns: 150,
	//	Rows:    35,
	//}

	//err, server := cmdbServerService.GetServerById(float64(sshClient.Server.ID))
	//if err != nil {
	//	conn.WriteMessage(1, []byte(err.Error()))
	//	conn.Close()
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//sshClient.Server = &server

	err = sshClient.GenerateClient()
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		conn.Close()
		response.FailWithMessage(err.Error(), c)
		return
	}
	//sshClient.RequestTerminal(terminal)
	//sshClient.Connect(conn)

	ssConn, err := common.NewSshConn(sshClient.Client, 150, 35)

	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		//conn.Close()
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer ssConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(conn, logBuff, quitChan)
	go ssConn.SendComboOutput(conn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
	return
}
