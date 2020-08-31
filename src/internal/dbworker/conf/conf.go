package conf

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	Conf *Config
	confPath string
	host string
)

func init() {
	host, _ = os.Hostname()
	flag.StringVar(&confPath, "conf", "comet.conf", "comet config path")
}

// Init init config.
func Init() (err error) {
	Conf = Default()

	_, err = toml.DecodeFile(confPath, &Conf)
	fmt.Println(Conf)
	return
}

// Default new a config with specified default value.
func Default() *Config {
	return &Config{
		MySql: &MySql{
			Username:  "",
			Password:  "",
			Hostname:  "",
			Port:      "",
			Database: "",
		},
		Kafka: &Kafka{
			Topic:   "test",
			Group:   "MySqlGroup",
			Brokers: []string{"localhost:9092"},
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
	Port string
	Database string

}
// Kafka is kafka config.
type Kafka struct {
	Topic   string
	Group   string
	Brokers []string
}
