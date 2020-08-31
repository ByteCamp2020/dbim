package comet

import (
	"bdim/src/api/comet/grpc"
	"bufio"
)

type Channel struct {
	Room   *Room
	signal chan *grpc.Package
	Next   *Channel
	Prev   *Channel
	Writer bufio.Reader
}

func NewChannel() *Channel {
	c := &Channel{
		Room:   nil,
		signal: make(chan *grpc.Package, 1024),
		Next:   nil,
		Prev:   nil,
		Writer: bufio.Reader{},
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

func (c *Channel) Listen() *grpc.Package {
	return <-c.signal
}
