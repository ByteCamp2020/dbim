package main

import (
	"bdim/src/internal/comet"
	"bdim/src/internal/comet/conf"
	"bdim/src/internal/comet/grpc"
	"bdim/src/models/discovery"
	"bdim/src/models/log"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Info("Initiating comet server.", nil)
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
	log.Print("Registering grpc server, add: " + cfg.RPCServer.RegAddr)
	d.RegComet(cfg.RPCServer.RegAddr)
	defer d.DelComet(cfg.RPCServer.RegAddr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-ch
		log.Info(fmt.Sprintf("bdim-comet get a signal %s", s.String()), nil)
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("bdim-comet exit", nil)
			grpcServer.GracefulStop()
			cm.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
