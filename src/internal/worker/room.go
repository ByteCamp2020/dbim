package worker

import (
	"bdim/src/pkg/bytes"
	"errors"
	"time"

	comet "bdim/src/api/comet/grpc"
	"bdim/src/internal/worker/conf"
	log "github.com/golang/glog"
)

var (
	// ErrRoomFull room chan full.
	ErrRoomFull = errors.New("room pack chan full")

	roomReadyPackage = new(comet.Package)
)

// Room room.
type Room struct {
	c     *conf.Room
	worker   *Worker
	id    int32
	pack chan *comet.Package
}

// NewRoom new a room struct, store channel room info.
func NewRoom(worker *Worker, id int32, c *conf.Room) (r *Room) {
	r = &Room{
		c:     c,
		id:    id,
		worker:   worker,
		pack: make(chan *comet.Package, c.Batch*2),
	}
	go r.pushproc(c.Batch, time.Duration(c.Signal))
	return
}

// Push push msg to the room, if chan full discard it.
func (r *Room) Push(op int32, msg []byte) (err error) {
	var p = &comet.Package{
		Op:   op,
		Body: msg,
	}
	select {
	case r.pack <- p:
	default:
		err = ErrRoomFull
	}
	return
}

// pushproc merge package and push msgs in batch.
func (r *Room) pushproc(batch int, sigTime time.Duration) {
	var (
		n    int
		last time.Time
		p    *comet.Package
		buf  = bytes.NewWriterSize(int(comet.MaxBodySize))
	)
	log.Infof("start room:%s goroutine", r.id)
	td := time.AfterFunc(sigTime, func() {
		select {
		case r.pack <- roomReadyPackage:
		default:
		}
	})
	defer td.Stop()
	for {
		if p = <-r.pack; p == nil {
			break // exit
		} else if p != roomReadyPackage {
			// merge buffer ignore error, always nil
			p.WriteTo(buf)
			if n++; n == 1 {
				last = time.Now()
				td.Reset(sigTime)
				continue
			} else if n < batch {
				if sigTime > time.Since(last) {
					continue
				}
			}
		} else {
			if n == 0 {
				break
			}
		}
		_ = r.worker.broadcastRoomRawBytes(r.id, buf.Buffer())
		// TODO use reset buffer
		// after push to room channel, renew a buffer, let old buffer gc
		buf = bytes.NewWriterSize(buf.Size())
		n = 0
		if r.c.Idle != 0 {
			td.Reset(time.Duration(r.c.Idle))
		} else {
			td.Reset(time.Minute)
		}
	}
	r.worker.delRoom(r.id)
	log.Infof("room:%s goroutine exit", r.id)
}


func (w *Worker) delRoom(roomID int32) {
	w.roomsMutex.Lock()
	delete(w.rooms, roomID)
	w.roomsMutex.Unlock()
}

func (w *Worker) getRoom(roomID int32) *Room {
	w.roomsMutex.RLock()
	room, ok := w.rooms[roomID]
	w.roomsMutex.RUnlock()
	if !ok {
		w.roomsMutex.Lock()
		if room, ok = w.rooms[roomID]; !ok {
			room = NewRoom(w, roomID, w.c.Room)
			w.rooms[roomID] = room
		}
		w.roomsMutex.Unlock()
		log.Infof("new a room:%s active:%d", roomID, len(w.rooms))
	}
	return room
}

