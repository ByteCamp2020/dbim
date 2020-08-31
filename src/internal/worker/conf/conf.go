package conf

import (
	xtime "bdim/src/pkg/time"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

var (
	confPath  string
	region    string
	zone      string
	deployEnv string
	host      string
	// Conf config
	Conf *Config
)

func init() {
	host, _ = os.Hostname()
	flag.StringVar(&confPath, "conf", "comet.conf", "comet config path")
}

// Init init config.
func Init() (err error) {
	Conf = Default()

	_, err = toml.DecodeFile(confPath, &Conf)
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
