package main

import (
	"bdim/src/internal/worker"
	"bdim/src/internal/worker/conf"
	"flag"
	log "github.com/golang/glog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	ver = "1.0"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Infof("bdim-worker start")

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
