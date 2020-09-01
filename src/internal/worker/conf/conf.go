package conf

import (
	xtime "bdim/src/pkg/time"
	"flag"
	"github.com/BurntSushi/toml"
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

	flag.StringVar(&confPath, "conf", "comet.conf", "comet config path")
	//flag.StringVar(&redisAddr, "redisAddr", "redis:localhost:6379", "redis")
	//flag.StringVar(&group, "kafkagroup", "bdim-worker", "kafka group")
	//flag.StringVar(&topic, "kafkatopic", "bdim", "kafka topic")
	//var broker string
	//flag.StringVar(&broker, "kafkabroker", "localhost:9092", "broker")
	//brokers = append(brokers, broker)
}

// Init init config.
func Init() (err error) {
	Conf = Default()
	//Conf.Discovery.RedisAddr = redisAddr
	//Conf.Kafka.Topic = topic
	//Conf.Kafka.Brokers = brokers
	//Conf.Kafka.Group = group
	_, err = toml.DecodeFile(confPath, &Conf)
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
