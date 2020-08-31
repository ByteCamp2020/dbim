package http

import (
	"bdim/src/internal/logic"
	"bdim/src/internal/logic/conf"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server is http server.
type Server struct {
	engine *gin.Engine
	logic  *logic.Logic
}

// New new a http server.
func New(c *conf.HTTPServer, l *logic.Logic) *Server {
	engine := gin.New()
	engine.Use(MWHandleErrors())
	go func() {
		if err := engine.Run(c.Addr); err != nil {
			panic(err)
		}
	}()
	s := &Server{
		engine: engine,
		logic:  l,
	}
	s.initRouter()
	return s
}

func (s *Server) initRouter() {
	group := s.engine.Group("/bdim")
	group.GET("/test", s.test)
	group.POST("/push", s.push)
}

// Close close the server.
func (s *Server) Close() {

}

type Arg struct {
	Op   int32  `form:"operation" binding:"required"`
	Room int32 `form:"room" binding:"required"`
	User string `form:"user" binding:"required"`
}

func (s *Server) test(c *gin.Context) {
	c.JSON(http.StatusOK, "hello")
}

func (s *Server) push(c *gin.Context) {
	var arg Arg
	if err := c.BindQuery(&arg); err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	// read message
	msg, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	if s.logic.DFA.CheckSentence(string(msg)) == false {
		return
	}
	{
		errors(c, RequestErr, err.Error())
		return
	}
	timestamp := int32(time.Now().Unix())
	if err = s.logic.PushRoom(c, arg.Op, arg.Room, arg.User, timestamp, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c, nil, OK)
}
