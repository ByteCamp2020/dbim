package conf

import (
	"flag"
	"os"
	"time"
)

var (
	// Conf config
	Conf *Config
	host string
	confPath string
)

func init() {
	host, _ = os.Hostname()
	flag.StringVar(&confPath, "conf", "comet.conf", "comet config path")
}

// Init init config.
func Init() (err error) {
	Conf = Default()
	return
}

// Default new a config with specified defualt value.
func Default() *Config {
	return &Config{
		HTTPServer: &HTTPServer{
			Network:      "tcp",
			Addr:         "localhost:2333",
			ReadTimeout:  time.Duration(time.Second),
			WriteTimeout: time.Duration(time.Second),
		},
		Kafka: &Kafka{
			Topic:   "test",
			Brokers: []string{"localhost:9092"},
		},
	}
}

// Config config.
type Config struct {
	HTTPServer *HTTPServer
	Kafka      *Kafka
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Kafka struct {
	Topic   string
	Brokers []string
}
