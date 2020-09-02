package worker

import (
	"context"
	"fmt"

	comet "bdim/src/api/comet/grpc"
	pb "bdim/src/api/logic/grpc"
	"bdim/src/models/log"
)

const opRaw = int32(0)

func (w *Worker) push(ctx context.Context, pushMsg *pb.PushMsg) (err error) {
	err = w.getRoom(pushMsg.Roomid).Push(pushMsg.Msg)
	return
}

// broadcastRoomRawBytes broadcast aggregation messages to room.
func (w *Worker) broadcastRoomRawBytes(roomID int32, body []byte) (err error) {
	args := comet.Package{
		Roomid: roomID,
		Body:   body,
	}
	comets := w.cometServers
	for serverID, c := range comets {
		fmt.Printf("c.BroadcastRoom(%v) roomID:%d serverID:%s\n", args.Body, roomID, serverID)
		if err = c.BroadcastRoom(&args); err != nil {
			log.Error(fmt.Sprintf("c.BroadcastRoom(%v) roomID:%v serverID:%s ", args.Body, roomID, serverID), err)
		}
	}
	log.Info(fmt.Sprintf("broadcastRoom comets:%d", len(comets)), nil)
	return
}
