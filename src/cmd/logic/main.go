package main

import (
	"bdim/src/internal/logic"
	"bdim/src/internal/logic/conf"
	"bdim/src/internal/logic/http"
	"bdim/src/models/log"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const (
	ver   = "0.0.1"
	appid = "bdim.logic"
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("bdim-logic [version: %s env: %+v] start", ver, conf.Conf), nil)

	// logic
	srv := logic.New(conf.Conf)

	httpSrv := http.New(conf.Conf.HTTPServer, srv)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c

		log.Info(fmt.Sprintf("bdim-logic get a signal %s", s.String()), nil)
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			httpSrv.Close()
			log.Info(fmt.Sprintf("bdim-logic [version: %s] exit", ver), nil)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
