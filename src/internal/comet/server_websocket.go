package comet

import (
	"bdim/src/internal/comet/conf"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var (
	registerCh = make(chan *Register)
)

type ClientManager struct {
	addr      string
	Clients   map[*Client]bool
	comet     *Comet
	cfg       *conf.WebSocket
	clientCnt int
}
type Client struct {
	channel *Channel
	conn    *websocket.Conn
}

type Register struct {
	conn   *websocket.Conn
	roomID int32
}

func NewClientManage(cfg *conf.WebSocket, comet *Comet) *ClientManager {
	cm := &ClientManager{
		addr:    cfg.WsAddr,
		Clients: make(map[*Client]bool, cfg.ClientNo),
		comet:   comet,
		cfg:     cfg,
	}
	go cm.registerPros()
	return cm
}

func NewClient(conn *websocket.Conn, c *Channel, roomID int32) *Client {
	client := &Client{
		channel: c,
		conn:    conn,
	}
	go client.pushProc()
	return client
}

func (c *Client) pushProc() {
	for {
		info := c.channel.Listen()
		fmt.Println(info)
		err := c.conn.WriteMessage(websocket.BinaryMessage, info.Body)
		if err != nil {
			return
		}
	}
}

func (cm *ClientManager) watch(c *Client) {
	for {
		_, _, err := c.conn.ReadMessage()
		fmt.Println("delected found")
		if err != nil {
			cm.del(c)
			return
		}
	}
}

func (cm *ClientManager) Close() {
	for k, _ := range cm.Clients {
		cm.del(k)
	}
}

func (cm *ClientManager) registerPros() {
	for {
		register := <-registerCh
		// new channel
		ch := NewChannel()
		cm.comet.Put(ch, register.roomID)
		client := NewClient(register.conn, ch, register.roomID)
		cm.Clients[client] = true
		go cm.watch(client)
	}
}

func (cm *ClientManager) del(c *Client) {
	c.channel.Room.Del(c.channel)
	c.conn.Close()
	delete(cm.Clients, c)
}

func StartWebSocket(addr string) {
	fmt.Println("start listening")
	http.HandleFunc("/push", serveHTTP)
	http.ListenAndServe(addr, nil)
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	fmt.Println("?????????????????")
	roomid, err := strconv.ParseInt(r.URL.Query()["roomid"][0], 10, 32)
	if err != nil {
		glog.Error("Args wrong", err)
	}
	conn, err := upgrade.Upgrade(w, r, nil)
	fmt.Println("??!!")
	if err != nil {
		fmt.Println(err)
		glog.Error("Upgrade fail", err)
	}
	register := &Register{
		conn:   conn,
		roomID: int32(roomid),
	}
	registerCh <- register
}
