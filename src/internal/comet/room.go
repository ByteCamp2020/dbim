package comet

import (
	"bdim/src/api/comet/grpc"
	"sync"
)

type Room struct {
	next   *Channel
	roomID int32
	rLock  sync.RWMutex
}

// NewRoom new a room struct, store channel room info.
func NewRoom(id int32) (r *Room) {
	r = new(Room)
	r.roomID = id
	r.next = nil
	return
}

func (r *Room) Put(ch *Channel) {
	r.rLock.Lock()
	if r.next != nil {
		r.next.Prev = ch
	}
	ch.Room = r
	ch.Next = r.next
	ch.Prev = nil
	r.next = ch
	r.rLock.Unlock()
}

func (r *Room) Push(p *grpc.Package) {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		_ = ch.Push(p)
	}
	r.rLock.RUnlock()
}

func (r *Room) Del(c *Channel) {
	r.rLock.Lock()
	if c.Next != nil {
		c.Next.Prev = c.Prev
	}
	if c.Prev != nil {
		c.Prev.Next = c.Next
	} else {
		r.next = c.Next
	}
	r.rLock.Unlock()
}

func (r *Room) Close() {
	var tmp *Channel
	var cur *Channel
	r.rLock.Lock()
	if r.next == nil {
		return
	}
	cur = r.next
	r.next = nil
	for cur != nil {
		tmp = cur.Next
		cur.Next = nil
		cur.Prev = nil
		cur = tmp
	}
	r.rLock.Unlock()
}
