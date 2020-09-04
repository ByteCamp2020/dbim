package conf

import (
	"log"
	"os"
)

var (
	Conf *Config
	host string
)

func init() {
	host, _ = os.Hostname()
	//flag.StringVar(&confPath, "conf", "comet.conf", "comet config path")
}

// Init init config.
func Init() (err error) {
	Conf = Default()

	//_, err = toml.DecodeFile(confPath, &Conf)
	log.Print(Conf)
	return
}

// Default new a config with specified default value.
func Default() *Config {
	return &Config{
		MySql: &MySql{
			Username: "wyb",
			Password: "123456",
			Hostname: "10.108.21.19",
			Port:     "3306",
			Database: "test",
		},
		Kafka: &Kafka{
			Topic:   "bdim",
			Group:   "MySqlGroup",
			Brokers: []string{"10.108.21.19:9092"},
		},
	}
}

type Config struct {
	Kafka *Kafka
	MySql *MySql
}

type MySql struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
}

// Kafka is kafka config.
type Kafka struct {
	Topic   string
	Group   string
	Brokers []string
}
