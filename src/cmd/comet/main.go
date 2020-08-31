package main

import (
	"bdim/src/internal/comet"
	"bdim/src/internal/comet/conf"
	"bdim/src/internal/comet/grpc"
	"bdim/src/models/discovery"
	"flag"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	glog.Info("Initiating comet server.")
	// new config
	cfg := conf.Init()
	// new comet
	c := comet.NewComet(cfg.Comet)

	// new room
	for i := 0; i < cfg.Comet.RoomNo; i++ {
		r := comet.NewRoom(int32(i))
		c.PutRoom(r)
	}
	// new client manager
	cm := comet.NewClientManage(cfg.WebSocket, c)
	// listen server go func()
	go comet.StartWebSocket(cfg.WebSocket.WsAddr)
	// new grpc server
	grpcServer := grpc.New(cfg.RPCServer, c)
	// register
	d := discovery.NewDiscovery(cfg.Discovery.RedisAddr)
	d.RegComet(cfg.RPCServer.Addr)
	defer d.DelComet(cfg.RPCServer.Addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-ch
		glog.Info("bdim-comet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			glog.Info("bdim-comet  exit")
			grpcServer.GracefulStop()
			cm.Close()
			glog.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
