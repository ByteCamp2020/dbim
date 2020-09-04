package comet

import (
	"bdim/src/internal/comet/conf"
	"bdim/src/models/log"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var (
	registerCh = make(chan *Register)
)

type ClientManager struct {
	addr      string
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
		info, ok := <- c.channel.signal
		if !ok {
			log.Print("STOP")
			return
		}
		log.Print(info)
		if c.conn == nil {
			return
		}
		err := c.conn.WriteMessage(websocket.BinaryMessage, info.Body)
		if err != nil {
			return
		}
	}
}

func (c *Client) watch() {
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Print("delected found,", err)
			c.del()
			return
		}
	}
}

func (cm *ClientManager) Close() {
	cm.comet.Close()
}

func (cm *ClientManager) registerPros() {
	for {
		register := <-registerCh
		// new channel
		ch := NewChannel()
		cm.comet.Put(ch, register.roomID)
		client := NewClient(register.conn, ch, register.roomID)
		register.conn = nil
		go client.watch()
	}
}

func (c *Client) del (){
	//log.Print("Before del")
	c.channel.Room.Del(c.channel)
	close(c.channel.signal)
	c.conn.Close()
	if c.conn != nil {
		c.conn = nil
	}
	//log.Print("After del")
}

func StartWebSocket(addr string) {
	log.Print("start listening")
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
	temprid := r.URL.Query()["roomid"]
	if (len(temprid) < 1) {
		log.Error("Args wrong", fmt.Errorf("roomid wrong"))
		return
	}

	roomid, err := strconv.ParseInt(temprid[0], 10, 32)
	if err != nil {
		log.Error("Args wrong", err)
		return
	}
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		log.Error("Upgrade fail", err)
		return
	}
	register := &Register{
		conn:   conn,
		roomID: int32(roomid),
	}
	log.Print(fmt.Sprintf("New connect to Room%v\n", roomid))
	registerCh <- register
}
