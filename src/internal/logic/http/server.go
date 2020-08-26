package http

import (
	"bdim/src/internal/logic"
	"bdim/src/internal/logic/conf"
	"net/http"

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
}

// Close close the server.
func (s *Server) Close() {

}

func (s *Server) test(c *gin.Context) {
	c.JSON(http.StatusOK, "hello")
}