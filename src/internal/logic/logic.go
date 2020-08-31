package logic

import (
	"context"

	pb "bdim/src/api/logic/grpc"
	"bdim/src/internal/logic/conf"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	kafka "gopkg.in/Shopify/sarama.v1"
)

// Logic struct
type Logic struct {
	c        *conf.Config
	kafkaPub kafka.SyncProducer
	DFA      *DFAUtil
}

// New init
func New(c *conf.Config) (l *Logic) {
	l = &Logic{
		c:        c,
		kafkaPub: newKafkaPub(c.Kafka),
		DFA:      NewDFAUtil(c.WordList),
	}
	return l
}

func newKafkaPub(c *conf.Kafka) kafka.SyncProducer {
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(c.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

// Close close resources.
func (l *Logic) Close() {

}

func (l *Logic) PushRoom(c context.Context, op int32, room int32, user string, timestamp int32, msg []byte) (err error) {
	pushMsg := &pb.PushMsg{
		Op:   op,
		Roomid: room,
		User: user,
		Msg:  msg,
	}
	Msg := &pb.Msg{
		Pm:        pushMsg,
		Timestamp: timestamp,
		Visible:   true,
	}

	b, err := proto.Marshal(Msg)
	if err != nil {
		return
	}
	m := &kafka.ProducerMessage{
		Key:   kafka.StringEncoder(room),
		Topic: l.c.Kafka.Topic,
		Value: kafka.ByteEncoder(b),
	}
	if _, _, err = l.kafkaPub.SendMessage(m); err != nil {
		log.Errorf("PushMsg.send(broadcast_room pushMsg:%v) error(%v)", pushMsg, err)
	}
	return
}
