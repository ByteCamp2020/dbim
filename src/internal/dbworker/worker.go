package dbworker

import (
	pb "bdim/src/api/logic/grpc"
	"bdim/src/internal/dbworker/conf"
	"fmt"

	"bdim/src/internal/dbworker/dao"
	cluster "github.com/bsm/sarama-cluster"
	"google.golang.org/protobuf/proto"
	log "github.com/golang/glog"
)

// Worker is push Worker.
type DbWorker struct {
	c        *conf.Config
	consumer *cluster.Consumer
	dao      *dao.Dao
}

// New new a push Worker.
func New(c *conf.Config) *DbWorker {
	w := &DbWorker{
		c:        c,
		consumer: newKafkaSub(c.Kafka),
		dao:      dao.New(c.MySql),
	}
	return w
}

func newKafkaSub(c *conf.Kafka) *cluster.Consumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(c.Brokers, c.Group, []string{c.Topic}, config)
	if err != nil {
		panic(err)
	}
	return consumer
}

// Close close resources.
func (w *DbWorker) Close() error {
	if w.dao != nil {
		err := w.dao.Close()
		if err != nil {
			fmt.Println("DbWorker: Close db err", err)
		}
	}
	if w.consumer != nil {
		return w.consumer.Close()
	}
	return nil
}

// Consume messages, watch signals
func (w *DbWorker) Consume() {
	for {
		select {
		case err := <-w.consumer.Errors():
			log.Errorf("consumer error(%v)", err)
		case n := <-w.consumer.Notifications():
			log.Infof("consumer rebalanced(%v)", n)
		case msg, ok := <-w.consumer.Messages():
			if !ok {
				return
			}
			w.consumer.MarkOffset(msg, "")
			// process push message
			mesg := new(pb.Msg)
			if err := proto.Unmarshal(msg.Value, mesg); err != nil {
				log.Errorf("proto.Unmarshal(%v) error(%v)", mesg, err)
				continue
			}
			fmt.Println("receive", msg)
			message := string(mesg.Pm.Msg)
			w.dao.AddMessage(mesg.Pm.User, mesg.Pm.Roomid, message, mesg.Timestamp, mesg.Visible)
			log.Infof("Dao:consume: %s/%d/%d\t%s\t%+v", msg.Topic, msg.Partition, msg.Offset, msg.Key, mesg)
		}
	}
}
