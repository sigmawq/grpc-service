package main

import (
	"context"
	"encoding/json"
	"errors"
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
	log.Printf("SendBatch: length %v", len(batch.Data))

	err := database.UpdateBatch(batch.Data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (server *Server) Retrieve(ctx context.Context, request *pb.RetrieveRequest) (*pb.RetrieveResponse, error) {
	log.Printf("Retrieve: from %v size %v search %v", request.From, request.Size, request.Search)

	data, err := database.Retrieve(request.Search, int(request.Size), int(request.From))
	if err != nil {
		return &pb.RetrieveResponse{}, errors.New("failed to retrieve data")
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return &pb.RetrieveResponse{}, errors.New("failed to retrieve data (serialization error)")
	}

	response := &pb.RetrieveResponse{
		Data:    string(serialized),
		Success: true,
	}
	return response, nil
}

func (server *Server) Aggregate(ctx context.Context, request *pb.AggregateRequest) (*pb.AggregateResponse, error) {
	log.Println("Aggregate")

	data, err := database.Aggregate()
	if err != nil {
		return &pb.AggregateResponse{}, errors.New("failed to retrieve data")
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return &pb.AggregateResponse{}, errors.New("failed to retrieve data (serialization error)")
	}

	response := &pb.AggregateResponse{
		Data:    string(serialized),
		Success: true,
	}
	return response, nil
}
