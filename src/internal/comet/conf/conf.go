package conf

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

var (
	confPath string
	host     string
	tp *TomlParas
)

type TomlParas struct {
	WsAddr string
	RPCAddr string
	RedisAddr string
}

type Config struct {
	WebSocket *WebSocket
	Comet     *Comet
	RPCServer *RPCServer
	Discovery *Discovery
}

func init() {
	host, _ = os.Hostname()
	//tp = &TomlParas{}
	flag.StringVar(&confPath, "conf", "comet.conf", "comet config path")
	//flag.StringVar(&tp.WsAddr, "wsaddr", ":3101", "websocket address")
	//flag.StringVar(&tp.RPCAddr, "rpcaddr", ":3109", "rpc server address")
	//flag.StringVar(&tp.RedisAddr, "redisaddr", ":6379", "redis addr")

}

func Default() *Config {
	return &Config{
		WebSocket: &WebSocket{
			ClientNo: int(1e7),
		},
		RPCServer: &RPCServer{
			Timeout:           time.Second,
			IdleTimeout:       time.Second * 60,
			MaxLifeTime:       time.Hour * 2,
			ForceCloseWait:    time.Second * 20,
			KeepAliveInterval: time.Second * 60,
			KeepAliveTimeout:  time.Second * 20,
		},
		Discovery: &Discovery{},
		Comet: &Comet{
			RoutinesNum: 8,
			RoutineSize: 1024,
			RoomNo:      1024,
			Host:        host,
		},
	}
}

func Init() *Config {
	cfg := Default()
	fmt.Println(confPath)

	if _, err := toml.DecodeFile(confPath, &tp); err != nil {
		fmt.Println(err)
	}
	fmt.Println(tp)
	//cfg.WebSocket.WsAddr = tp.WsAddr
	//cfg.RPCServer.Addr = tp.RPCAddr
	//cfg.Discovery.RedisAddr = tp.RedisAddr
	return cfg
}

// redis addr
// grpc addr
// port
type RPCServer struct {
	Addr              string
	Timeout           time.Duration
	IdleTimeout       time.Duration
	MaxLifeTime       time.Duration
	ForceCloseWait    time.Duration
	KeepAliveInterval time.Duration
	KeepAliveTimeout  time.Duration
}

type Comet struct {
	RoutinesNum uint64
	RoutineSize uint64
	RoomNo      int
	Host        string
}

type Discovery struct {
	RedisAddr string
}

type WebSocket struct {
	WsAddr   string
	ClientNo int
}
