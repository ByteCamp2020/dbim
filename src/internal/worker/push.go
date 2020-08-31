package worker

import (
	"context"
	"fmt"

	comet "bdim/src/api/comet/grpc"
	pb "bdim/src/api/logic/grpc"

	log "github.com/golang/glog"
)

const opRaw = int32(0)

func (w *Worker) push(ctx context.Context, pushMsg *pb.PushMsg) (err error) {
	err = w.getRoom(pushMsg.Roomid).Push(pushMsg.Op, pushMsg.Msg)
	return
}

// broadcastRoomRawBytes broadcast aggregation messages to room.
func (w *Worker) broadcastRoomRawBytes(roomID int32, body []byte) (err error) {
	args := comet.Package{
		Op:     opRaw,
		Roomid: roomID,
		Body:   body,
	}
	comets := w.cometServers
	for serverID, c := range comets {
		fmt.Printf("c.BroadcastRoom(%v) roomID:%d serverID:%s\n", args, roomID, serverID)
		if err = c.BroadcastRoom(&args); err != nil {
			log.Errorf("c.BroadcastRoom(%v) roomID:%s serverID:%s error(%v)", args, roomID, serverID, err)
		}
	}
	log.Infof("broadcastRoom comets:%d", len(comets))
	return
}
