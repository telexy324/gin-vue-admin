package sockets

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type connection struct {
	ws     *websocket.Conn
	send   chan []byte
	userID int
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		global.GVA_LOG.Error(c.ws.Close().Error(), zap.Any("error", "Error closing websocket"))
	}()

	c.ws.SetReadLimit(maxMessageSize)
	global.GVA_LOG.Error(c.ws.SetReadDeadline(time.Now().Add(pongWait)).Error(), zap.Any("error", "Socket state corrupt"))
	c.ws.SetPongHandler(func(string) error {
		global.GVA_LOG.Error(c.ws.SetReadDeadline(time.Now().Add(pongWait)).Error(), zap.Any("error", "Socket state corrupt"))
		return nil
	})

	for {
		_, message, e := c.ws.ReadMessage()
		fmt.Println(string(message))

		if e != nil {
			if websocket.IsUnexpectedCloseError(e, websocket.CloseGoingAway) {
				global.GVA_LOG.Error(e.Error())
			}
			break
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	global.GVA_LOG.Error(c.ws.SetWriteDeadline(time.Now().Add(writeWait)).Error(), zap.Any("error", "Socket state corrupt"))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		global.GVA_LOG.Error(c.ws.Close().Error())
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.write(websocket.CloseMessage, []byte{})
				if err != nil {
					global.GVA_LOG.Error(err.Error())
				}
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				global.GVA_LOG.Error(err.Error())
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				global.GVA_LOG.Error(err.Error())
				return
			}
		}
	}
}

// Handler is used by the router to handle the /ws endpoint
func Handler(c *gin.Context) {
	user := context.Get(c.Request, "user").(*system.SysUser)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	conn := &connection{
		send:   make(chan []byte, 256),
		ws:     ws,
		userID: int(user.ID),
	}

	h.register <- conn

	go conn.writePump()
	conn.readPump()
}

// Message allows a message to be sent to the websockets, called in API task logging
func Message(userID int, message []byte) {
	h.broadcast <- &sendRequest{
		userID: userID,
		msg:    message,
	}
}
