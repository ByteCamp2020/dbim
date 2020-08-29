package worker

import (
	pb "bdim/src/api/logic/grpc"
	"bdim/src/internal/worker/conf"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bilibili/discovery/naming"
	"github.com/gogo/protobuf/proto"

	cluster "github.com/bsm/sarama-cluster"
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
	}
	w.watchComet(c.Discovery)
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
			pushMsg := new(pb.PushMsg)
			if err := proto.Unmarshal(msg.Value, pushMsg); err != nil {
				log.Errorf("proto.Unmarshal(%v) error(%v)", msg, err)
				continue
			}
			if err := w.push(context.Background(), pushMsg); err != nil {
				log.Errorf("w.push(%v) error(%v)", pushMsg, err)
			}
			log.Infof("consume: %s/%d/%d\t%s\t%+v", msg.Topic, msg.Partition, msg.Offset, msg.Key, pushMsg)
		}
	}
}

func (w *Worker) watchComet(c *naming.Config) {
	dis := naming.New(c)
	resolver := dis.Build("bdim.comet")
	event := resolver.Watch()
	select {
	case _, ok := <-event:
		if !ok {
			panic("watchComet init failed")
		}
		if ins, ok := resolver.Fetch(); ok {
			if err := w.newAddress(ins.Instances); err != nil {
				panic(err)
			}
			log.Infof("watchComet init newAddress:%+v", ins)
		}
	case <-time.After(10 * time.Second):
		log.Error("watchComet init instances timeout")
	}
	go func() {
		for {
			if _, ok := <-event; !ok {
				log.Info("watchComet exit")
				return
			}
			ins, ok := resolver.Fetch()
			if ok {
				if err := w.newAddress(ins.Instances); err != nil {
					log.Errorf("watchComet newAddress(%+v) error(%+v)", ins, err)
					continue
				}
				log.Infof("watchComet change newAddress:%+v", ins)
			}
		}
	}()
}

func (w *Worker) newAddress(insMap map[string][]*naming.Instance) error {
	ins := insMap[w.c.Env.Zone]
	if len(ins) == 0 {
		return fmt.Errorf("watchComet instance is empty")
	}
	comets := map[string]*Comet{}
	for _, in := range ins {
		if old, ok := w.cometServers[in.Hostname]; ok {
			comets[in.Hostname] = old
			continue
		}
		c, err := NewComet(in, w.c.Comet)
		if err != nil {
			log.Errorf("watchComet NewComet(%+v) error(%v)", in, err)
			return err
		}
		comets[in.Hostname] = c
		log.Infof("watchComet AddComet grpc:%+v", in)
	}
	for key, old := range w.cometServers {
		if _, ok := comets[key]; !ok {
			old.cancel()
			log.Infof("watchComet DelComet:%s", key)
		}
	}
	w.cometServers = comets
	return nil
}
