package comet

import (
	"bdim/src/api/comet/grpc"
	"bdim/src/internal/comet/conf"
	"sync/atomic"
)

type Comet struct {
	serverID string
	rooms map[int32]*Room
	routines []chan *grpc.Package

	routinesNum uint64
	routineAmount uint64
}

func NewComet(cfg *conf.Config) *Comet{
	c := &Comet {
		serverID: cfg.Host,
		routinesNum: 0,
		routineAmount: cfg.RoutinesNum,
	}
	for i := uint64(0); i < c.routineAmount; i++ {
		ch := make(chan *grpc.Package, cfg.RoutineSize)
		c.routines[i] = ch
		go c.cometProc(ch)
	}
	return c
}
func (c *Comet) Put (ch *Channel, roomID int32) {
	c.Room(roomID).Put(ch)
}
func (c *Comet) Room(roomID int32) (room *Room){
	room = c.rooms[roomID]
	return
}

func (c *Comet) Push(p *grpc.Package) {
	idx := atomic.AddUint64(&c.routinesNum, 1) % c.routineAmount
	c.routines[idx] <- p
}

func (c *Comet) cometProc(ch chan *grpc.Package) {
	for {
		pack := <-ch
		if room := c.Room(pack.Roomid); room != nil {
			room.Push(pack)
		}
	}
}