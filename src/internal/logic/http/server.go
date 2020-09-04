package http

import (
	"bdim/src/internal/logic"
	"bdim/src/internal/logic/conf"
	"bdim/src/models/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
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
	group.GET("query", s.query)
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
	log.Print( "  here is a push request !!!!!!!:w")
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
	if (s.logic.C.HTTPServer.IsForbidden == 1 && s.logic.DFA != nil) && s.logic.DFA.CheckSentence(string(msg)) == false {
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
	result(c, OK)
}

func (s *Server) query(c *gin.Context) {
	uid := c.Query("uid")
	roomid := c.Query("roomid")
	timestamp := c.Query("timestamp")
	if roomid != "" {
		_, err := strconv.Atoi(roomid)
		if err != nil {
			errors(c, RequestErr, err.Error())
			return
		}
	}
	if timestamp != "" {
		_, err := strconv.Atoi(timestamp)
		if err != nil {
			errors(c, RequestErr, err.Error())
			return
		}
	}
	res, err := s.logic.DbC.GetMessage(uid, roomid, timestamp)
	if err != nil {
		errors(c, ServerErr, err.Error())
		return
	}
	c.Set(contextErrCode, OK)
	c.JSON(200, res)
}
