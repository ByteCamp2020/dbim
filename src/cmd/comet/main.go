package comet

import (
	"bdim/src/internal/comet/conf"
	"bdim/src/internal/comet"
	"bdim/src/internal/comet/grpc"

)
func main() {
	// new config
	cfg := conf.Init()
	// new comet
	c := comet.NewComet(cfg)
	// register service
	// new room

	// new client manager
	cm := comet.NewClientManage(cfg, c)
	// new grpc server
	grpcServer := grpc.New(cfg, c)
	// register
	// listen server go func()
	go comet.StartWebSocket()
}