package main

import (
	"bdim/src/internal/worker"
	"bdim/src/internal/worker/conf"
	"bdim/src/models/log"
	"flag"
	"fmt"
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
	log.Info("bdim-worker start", nil)

	// worker
	w := worker.New(conf.Conf)
	go w.Consume()
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		log.Info(fmt.Sprintf("bdim-worker get a signal %s", s.String()),nil)
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			w.Close()
			log.Info(fmt.Sprintf("bdim-worker [version: %s] exit", ver),nil)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
