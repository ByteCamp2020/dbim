package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"bdim/src/internal/logic"
	"bdim/src/internal/logic/conf"
	"bdim/src/internal/logic/http"
	log "github.com/golang/glog"
)

const (
	ver   = "0.0.1"
	appid = "dbim.logic"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Infof("dbim-logic [version: %s env: %+v] start", ver, conf.Conf)

	// logic
	srv := logic.New(conf.Conf)
	httpSrv := http.New(conf.Conf.HTTPServer, srv)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("dbim-logic get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			httpSrv.Close()
			log.Infof("dbim-logic [version: %s] exit", ver)
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
