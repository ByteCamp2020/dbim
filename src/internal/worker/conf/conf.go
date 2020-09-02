package conf

import (
	xtime "bdim/src/pkg/time"
	"fmt"
	"os"
	"time"
)

var (
	confPath  string
	redisAddr string
	group string
	topic string
	brokers []string
	host string
	// Conf config
	Conf *Config
)

func init() {
	host, _ = os.Hostname()
	redisAddr = "redis://10.108.21.18:6379"
	group = "bdim-worker"
	topic = "kafkatopic"
	var broker string
	broker ="10.108.21.19:9092"
	brokers = append(brokers, broker)
}

// Init init config.
func Init() (err error) {
	Conf = Default()
	Conf.Discovery.RedisAddr = redisAddr
	Conf.Kafka.Topic = topic
	Conf.Kafka.Brokers = brokers
	Conf.Kafka.Group = group
	fmt.Println(Conf.Kafka)
	return
}

// Default new a config with specified defualt value.
func Default() *Config {
	return &Config{
		Discovery: &Discovery{RedisAddr: ":6379"},
		Comet:     &Comet{RoutineChan: 1024, RoutineSize: 32},
		Room: &Room{
			Batch:  20,
			Signal: xtime.Duration(time.Second),
			Idle:   xtime.Duration(time.Minute * 15),
		},
		Kafka: &Kafka{
			Topic:   "test",
			Brokers: []string{"localhost:9092"},
		},
	}
}

// Config is worker config.
type Config struct {
	Kafka     *Kafka
	Discovery *Discovery
	Comet     *Comet
	Room      *Room
}

// Room is room config.
type Room struct {
	Batch  int
	Signal xtime.Duration
	Idle   xtime.Duration
}

// Comet is comet config.
type Comet struct {
	RoutineChan int
	RoutineSize int
}

// Kafka is kafka config.
type Kafka struct {
	Topic   string
	Group   string
	Brokers []string
}

type Discovery struct {
	RedisAddr string
}
