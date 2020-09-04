package main

import (
	"bdim/src/internal/dbworker"
	"bdim/src/internal/dbworker/conf"
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
	if err := conf.Init(); err != nil {
		panic(err)
	}
	c := conf.Conf
	log.Info(fmt.Sprintf("DbWorker: Starting dbworker, cfg:%v", c), nil)
	dbWorker := dbworker.New(c)
	go dbWorker.Consume()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-wait

	err := dbWorker.Close()
	if err != nil {
		log.Error("Failed to close consumer: %v", err)
	}
}
