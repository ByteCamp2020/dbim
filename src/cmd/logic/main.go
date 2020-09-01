package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bdim/src/internal/logic"
	"bdim/src/internal/logic/conf"
	"bdim/src/internal/logic/http"
	log "github.com/golang/glog"
)

const (
	ver   = "0.0.1"
	appid = "bdim.logic"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Infof("bdim-logic [version: %s env: %+v] start", ver, conf.Conf)

	// logic
	srv := logic.New(conf.Conf)
	conf.Conf.HTTPServer.IsLimit = true
	conf.Conf.HTTPServer.RedisAddr = "redis://localhost:6379"
	// 2 times in 1 second
	conf.Conf.HTTPServer.Count = 2
	conf.Conf.HTTPServer.Dur = 1 * time.Second

	httpSrv := http.New(conf.Conf.HTTPServer, srv)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("bdim-logic get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			httpSrv.Close()
			log.Infof("bdim-logic [version: %s] exit", ver)
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
