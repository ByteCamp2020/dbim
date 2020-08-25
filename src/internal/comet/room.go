package comet

import (
	"bdim/src/api/comet/grpc"
	"sync"
)

type Room struct {
	next *Channel
	roomID string
	rLock sync.RWMutex
}

func (r *Room) Put(ch *Channel) {
	r.rLock.Lock()
	if r.next != nil {
		r.next.Prev = ch
	}
	ch.Next = r.next
	ch.Prev = nil
	r.next = ch
	r.rLock.Unlock()
}

func (r *Room) Push(p *grpc.Package){
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		_ = ch.Push(p)
	}
	r.rLock.RUnlock()
}