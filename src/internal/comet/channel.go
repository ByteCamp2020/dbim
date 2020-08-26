package comet

import (
	"bdim/src/api/comet/grpc"
	"bdim/src/internal/comet/conf"
	"bufio"
)
type Channel struct {
	Room *Room
	signal chan *grpc.Package
	Next *Channel
	Prev *Channel
	Writer bufio.Reader
}

func NewChannel(cfg *conf.Config) *Channel {
	c := &Channel{
		Room:   nil,
		signal: make(chan *grpc.Package, cfg.RoutineSize),
		Next:   nil,
		Prev:   nil,
		Writer: bufio.Reader{},
	}

	return c
}

func (c *Channel) Push(p *grpc.Package) (err error){
	select {
	case c.signal <- p:
	default:
	}
	return
}

func (c *Channel) Listen() *grpc.Package {
	return <- c.signal
}