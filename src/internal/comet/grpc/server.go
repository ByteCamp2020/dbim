package grpc

import (
	pb "bdim/src/api/comet/grpc"
	"bdim/src/internal/comet"
	"bdim/src/internal/comet/conf"
	"context"
	"flag"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	c *comet.Comet
}

func New(cfg *conf.Config, c *comet.Comet) *grpc.Server {
	flag.Parse()
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {

	}
	grpcServer := grpc.NewServer()
	pb.RegisterCometServer(grpcServer, &server{c})
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return grpcServer
}

func (s *server) Push(ctx context.Context, p *pb.Package) (*pb.PushReply, error){
	s.c.Push(p)
	resp := &pb.PushReply{}
	return resp, nil
}