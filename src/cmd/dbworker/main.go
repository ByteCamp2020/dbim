package main

import (
	"bdim/src/internal/dbworker"
	"bdim/src/internal/dbworker/conf"
	"flag"
	log "github.com/golang/glog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := conf.Init(); err != nil {
		panic(err)
	}
	c := conf.Conf
	log.Infof("DbWorker: Starting dbworker, cfg:%v", c)
	dbWorker := dbworker.New(c)
	go dbWorker.Consume()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-wait

	err := dbWorker.Close()
	if err != nil {
		log.Errorf("Failed to close consumer: %v", err)
	}
}
