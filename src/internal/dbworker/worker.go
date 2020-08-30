package dbworker

import (
	pb "bdim/src/api/logic/grpc"
	"bdim/src/internal/dbworker/conf"
	"context"
	"fmt"

	"bdim/src/models/discovery"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	"bdim/src/internal/dbworker/dao"
)

// Worker is push Worker.
type DbWorker struct {
	c            *conf.Config
	consumer     *cluster.Consumer
	dao *dao.Dao
}

// New new a push Worker.
func New(c *conf.Config) *DbWorker {
	w := &DbWorker{
		c:        c,
		consumer: newKafkaSub(c.Kafka),
		dao: dao.New(c.MySql),
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
			pushMsg := new(pb.PushMsg)
			if err := proto.Unmarshal(msg.Value, pushMsg); err != nil {
				log.Errorf("proto.Unmarshal(%v) error(%v)", msg, err)
				continue
			}
			fmt.Println("receive", pushMsg)
			w.dao.AddMessage(pushMsg.User, pushMsg.Roomid, pushMsg.)
			log.Infof("Dao:consume: %s/%d/%d\t%s\t%+v", msg.Topic, msg.Partition, msg.Offset, msg.Key, pushMsg)
		}
	}
}

func (w *Worker) initComet(c *conf.Discovery) {
	dis := discovery.NewDiscovery(c.RedisAddr)
	cometAddrs := dis.GetCometAddr()
	fmt.Println(cometAddrs)
	for _, addr := range cometAddrs {
		cmt, _ := NewComet(addr, w.c.Comet)
		w.cometServers[addr] = cmt
	}
}

