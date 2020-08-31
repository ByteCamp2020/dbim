package worker

import (
	pb "bdim/src/api/logic/grpc"
	"bdim/src/internal/worker/conf"
	"context"
	"fmt"
	"sync"

	"bdim/src/models/discovery"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
)

// Worker is push Worker.
type Worker struct {
	c            *conf.Config
	consumer     *cluster.Consumer
	cometServers map[string]*Comet

	rooms      map[int32]*Room
	roomsMutex sync.RWMutex
}

// New new a push Worker.
func New(c *conf.Config) *Worker {
	w := &Worker{
		c:        c,
		consumer: newKafkaSub(c.Kafka),
		rooms:    make(map[int32]*Room),
		cometServers: make(map[string]*Comet),
	}
	w.initComet(c.Discovery)
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
func (w *Worker) Close() error {
	if w.consumer != nil {
		return w.consumer.Close()
	}
	return nil
}

// Consume messages, watch signals
func (w *Worker) Consume() {
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
				log.Errorf("proto.Unmarshal(%v) error(%v)", msg, err)
				continue
			}
			pushMsg := mesg.Pm
			fmt.Println("receive", pushMsg)
			if err := w.push(context.Background(), pushMsg); err != nil {
				log.Errorf("w.push(%v) error(%v)", pushMsg, err)
			}
			log.Infof("consume: %s/%d/%d\t%s\t%+v", msg.Topic, msg.Partition, msg.Offset, msg.Key, pushMsg)
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

