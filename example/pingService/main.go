package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "pingService/proto"
)

var (
	port = flag.Int("port", 9600, "The server port")
)

type pingServer struct {
	pb.UnimplementedPingServiceServer
}

func newServer() *pingServer {
	return &pingServer{}
}

func (p *pingServer) Ping(context context.Context, pingRequest *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Timestamp: 1, Message: "test"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPingServiceServer(grpcServer, newServer())
	log.Printf("Started service at port %d...", *port)
	grpcServer.Serve(lis)
	log.Printf("Ping service stopped")
}
