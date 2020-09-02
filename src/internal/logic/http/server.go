package http

import (
	"bdim/src/internal/logic"
	"bdim/src/internal/logic/conf"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

// Server is http server.
type Server struct {
	engine *gin.Engine
	logic  *logic.Logic
}

// New new a http server.
func New(c *conf.HTTPServer, l *logic.Logic) *Server {
	engine := gin.New()

	if c.IsLimit == true {
		limiter, err := logic.NewLimiter(*c)
		if err != nil {
			panic(err)
		}
		engine.Use(MWHandleErrors(), RateMiddleware(limiter))
		l.Limiter = limiter
	} else {
		engine.Use(MWHandleErrors())
	}
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
	Room int32  `form:"room" binding:"required"`
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
	// check the forbidden words
	if s.logic.DFA.CheckSentence(string(msg)) == false {
		errors(c, RequestErr, "forbidden word")
		return
	}
	// user limit
	timestamp := int32(time.Now().Unix())
	if s.logic.Limiter != nil && s.logic.Limiter.UserLimit(arg.User) == false {
		errors(c, http.StatusTooManyRequests, "too many requests from the same user")
		return
	}

	if err = s.logic.PushRoom(c, arg.Room, arg.User, timestamp, msg); err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	result(c,  OK)
}
