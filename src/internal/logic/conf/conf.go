package conf

import (
	"time"
)

var (
	host      string

	// Conf config
	Conf *Config
)

func init() {
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
	}
}

// Config config.
type Config struct {
	HTTPServer *HTTPServer
}

// HTTPServer is http server config.
type HTTPServer struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
