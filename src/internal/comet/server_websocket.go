package comet

import (
	"bdim/src/internal/comet/conf"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var (
	registerCh = make(chan *Register)
)

type ClientManager struct {
	addr string
	Clients map[*Client]bool
	comet *Comet
	cfg *conf.Config
}
type Client struct {
	channel *Channel
	conn *websocket.Conn
}

type Register struct {
	conn *websocket.Conn
	roomID int32
}

func NewClientManage(cfg *conf.Config, comet *Comet) *ClientManager {
	cm := &ClientManager{
		addr:    cfg.WsAddr,
		Clients: make(map[*Client]bool, cfg.Client),
		comet:   comet,
		cfg:     cfg,
	}
	cm.registerPros()
	return cm
}

func NewClient(conn *websocket.Conn, c *Channel, roomID int32) *Client{
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
		register := <-registerCh
		// new channel
		ch := NewChannel(cm.cfg)
		cm.comet.Put(ch, register.roomID)
		client := NewClient(register.conn, ch, register.roomID)
		cm.Clients[client] = true
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
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	roomid, err := strconv.ParseInt(r.Header["roomid"][0], 10, 32)
	if err != nil {

	}
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {

	}
	register := &Register{
		conn:    conn,
		roomID:  int32(roomid),
	}
	registerCh <- register
}
