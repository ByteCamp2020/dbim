package comet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"time"
)

var (
	registerCh = make(chan *Client)
)
type ClientManager struct {
	addr string
	upgrade *websocket.Upgrader
	Clients []*Client
}
type Client struct {
	channel *Channel
	conn *websocket.Conn
}

func NewClientManage() *ClientManager {
	return nil
}

func NewClient(conn *websocket.Conn, c *Channel) *Client{
	client := &Client{
		channel: c,
		conn: conn,
	}
	go client.pushProc()
	return client
}

func (c *Client) pushProc (){
	for {
		info := <-c.channel.signal
		c.conn.WriteMessage(websocket.BinaryMessage, info.Body)
	}
}

func (cm *ClientManager) registerPros() {
	for {
		conn := <-registerCh
		// new channel

		client := NewClient(conn, nil) //TODO:add channel
		cm.Clients = append(cm.Clients, client)
	}
}

func StartWebSocket() {
	http.HandleFunc("/test", ServeHTTP)
	http.ListenAndServe(":8089", nil)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}
	roomid := r.Header["roomid"]

	conn, err := upgrade.Upgrade(w, r, nil)
	// get room id
	// TODO:封装conn，在comet.room中添加channel
	if err != nil {

	}
	registerCh <- conn
}
