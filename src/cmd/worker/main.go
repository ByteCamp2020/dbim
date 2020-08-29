package main

import (
	"bdim/src/internal/worker"
	"bdim/src/internal/worker/conf"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bilibili/discovery/naming"

	resolver "github.com/bilibili/discovery/naming/grpc"
	log "github.com/golang/glog"
)

var (
	ver = "1.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Infof("bdim-work [version: %s env: %+v] start", ver, conf.Conf.Env)

	// grpc register naming
	dis := naming.New(conf.Conf.Discovery)
	resolver.Register(dis)

	// worker
	w := worker.New(conf.Conf)
	go w.Consume()

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		log.Infof("bdim-worker get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			w.Close()
			log.Infof("bdim-worker [version: %s] exit", ver)
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
