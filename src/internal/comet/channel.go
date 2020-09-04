package comet

import (
	"bdim/src/api/comet/grpc"
)

type Channel struct {
	Room   *Room
	signal chan *grpc.Package
	Next   *Channel
	Prev   *Channel
}

func NewChannel() *Channel {
	c := &Channel{
		Room:   nil,
		signal: make(chan *grpc.Package, 1024),
		Next:   nil,
		Prev:   nil,
	}
	return c
}

func (c *Channel) Push(p *grpc.Package) (err error) {
	select {
	case c.signal <- p:
	default:
	}
	return
}

//func (c *Channel) Listen() (*grpc.Package, error) {
//	return <-c.signal
//}
