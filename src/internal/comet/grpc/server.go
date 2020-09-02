package grpc

import (
	pb "bdim/src/api/comet/grpc"
	"bdim/src/internal/comet"
	"bdim/src/internal/comet/conf"
	"bdim/src/models/log"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
)

type server struct {
	c *comet.Comet
}

func New(cfg *conf.RPCServer, c *comet.Comet) *grpc.Server {
	keepParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.IdleTimeout,
		MaxConnectionAgeGrace: cfg.ForceCloseWait,
		Time:                  cfg.KeepAliveInterval,
		Timeout:               cfg.KeepAliveTimeout,
		MaxConnectionAge:      cfg.MaxLifeTime,
	})
	log.Print(cfg.Addr)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s", cfg.Addr))
	if err != nil {

	}
	grpcServer := grpc.NewServer(keepParams)
	pb.RegisterCometServer(grpcServer, &server{c})
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Print(err)
			panic(err)
		}
	}()
	return grpcServer
}

func (s *server) Push(ctx context.Context, p *pb.Package) (*pb.PushReply, error) {
	s.c.Push(p)
	log.Print("Receive message", p.Body)
	resp := &pb.PushReply{}
	return resp, nil
}

func (s *server) Close() {

}
