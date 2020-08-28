package conf

import (
	"os"
	"time"
)

var (
	host, _ = os.Hostname()
)

type Config struct {
	WebSocket *WebSocket
	Comet *Comet
	RPCServer *RPCServer
	Discovery *Discovery
}

func Init() *Config {
	return &Config{

		WebSocket: &WebSocket{
			WsAddr: ":3101",
			ClientNo: int(1e7),
		},
		RPCServer: &RPCServer{
			Addr:              ":3109",
			Timeout:           time.Second,
			IdleTimeout:       time.Second * 60,
			MaxLifeTime:       time.Hour * 2,
			ForceCloseWait:    time.Second * 20,
			KeepAliveInterval: time.Second * 60,
			KeepAliveTimeout:  time.Second * 20,
		},
		Discovery: &Discovery{
			RedisAddr: ":6379",
		},
		Comet: &Comet{
			RoutinesNum: 8,
			RoutineSize: 1024,
			RoomNo:      1024,
			Host: host,
		},
	}
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
	RoomNo int
	Host string
}

type Discovery struct {
	RedisAddr string
}

type WebSocket struct {
	WsAddr string
	ClientNo int
}
