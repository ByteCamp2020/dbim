package logic

import (
	pb "bdim/src/api/logic/grpc"
	"bdim/src/internal/logic/conf"
	"bdim/src/models/log"
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	kafka "gopkg.in/Shopify/sarama.v1"
)

// Logic struct
type Logic struct {
	C        *conf.Config
	kafkaPub kafka.SyncProducer
	DFA      *DFAUtil
	Limiter  *Limiter
	DbC      *dbc
}

// New init
func New(c *conf.Config) (l *Logic) {
	l = &Logic{
		C:        c,
		kafkaPub: newKafkaPub(c.Kafka),
		DbC:      NewDbC(c.MySql),
	}
	if c.HTTPServer.IsForbidden == 1 {
		l.DFA = NewDFAUtil(c.WordList)
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

func (l *Logic) PushRoom(c context.Context, room int32, user string, timestamp int32, msg []byte) (err error) {
	pushMsg := &pb.PushMsg{
		Roomid: room,
		User:   user,
		Msg:    msg,
	}
	Msg := &pb.Msg{
		Pm:        pushMsg,
		Timestamp: timestamp,
		Visible:   true,
	}
	b, err := proto.Marshal(Msg)
	if err != nil {
		return err
	}
	m := &kafka.ProducerMessage{
		Key:   kafka.StringEncoder(room),
		Topic: l.C.Kafka.Topic,
		Value: kafka.ByteEncoder(b),
	}
	if _, _, err = l.kafkaPub.SendMessage(m); err != nil {
		log.Error(fmt.Sprintf("PushMsg.send(broadcast_room pushMsg:%v) error ", pushMsg), err)
		return err
	}
	return
}
