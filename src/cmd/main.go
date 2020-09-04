package main

import (
	"bdim/src/internal/logic"
	"bdim/src/models/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Server is http server.
type Server struct {
	engine *gin.Engine
	logic  *logic.Logic
}

func main() {
	addr := "localhost:2333"
	StartWebSocket(addr)
}

func StartWebSocket(addr string) {
	log.Print("start listening")
	http.HandleFunc("/bdim/push", serveHTTP)
	http.ListenAndServe(addr, nil)
}
// curl -d 'hello' "http://localhost:2333/bdim/push?room=1&user=1"
func serveHTTP(w http.ResponseWriter, r *http.Request) {
	rid := r.URL.Query()["room"]
	usrid := r.URL.Query()["user"]

	roomid, err := strconv.ParseInt(rid[0], 10, 32)
	userid, err := strconv.ParseInt(usrid[0], 10, 32)
	if err != nil {
		log.Error("Args wrong", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}
	fmt.Println(body)
	fmt.Println(roomid)
	fmt.Println(userid)
}



