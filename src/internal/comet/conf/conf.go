package conf

import (
	"bdim/src/models/log"
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
	RPCRegAddr string
}

type Config struct {
	WebSocket *WebSocket
	Comet     *Comet
	RPCServer *RPCServer
	Discovery *Discovery
	Host string
}

func init() {
	host, _ = os.Hostname()
	tp = &TomlParas{}
	tp.WsAddr ="0.0.0.0:3101"
	tp.RPCAddr = "0.0.0.0:3209"
	tp.RPCRegAddr = os.Getenv("COMET_GRPC_SERVER")
	tp.RedisAddr = "redis://10.108.21.18:6379"
}

func Default() *Config {
	return &Config{
		WebSocket: &WebSocket{
			RoomNo: int(1024),
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
		Host: host,
	}
}

func Init() *Config {
	cfg := Default()
	cfg.WebSocket.WsAddr = tp.WsAddr
	cfg.RPCServer.Addr = tp.RPCAddr
	cfg.RPCServer.RegAddr = tp.RPCRegAddr
	cfg.Discovery.RedisAddr = tp.RedisAddr
	log.Print(tp)
	return cfg
}

// redis addr
// grpc addr
// port
type RPCServer struct {
	Addr              string
	RegAddr 		  string
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
	RoomNo int
}
