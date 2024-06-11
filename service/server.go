package main

import (
	"context"
	pb "github.com/sigmawq/grpc-service/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.UnimplementedGrpcServer
}

func NewServer(host string) (Server, error) {
	server := Server{}

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Printf("Failed to begin listen on %v: %v", host, err)
		return server, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGrpcServer(grpcServer, &server)

	log.Printf("Listen on %v", host)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("Failed to initialize GRPC server: %v", err)
		return server, err
	}

	return server, nil
}

func (server *Server) SendBatch(ctx context.Context, batch *pb.Batch) (*pb.BatchResponse, error) {
	log.Printf("Receive batch length %v: %#v", len(batch.Data), batch)
	return nil, nil
}
