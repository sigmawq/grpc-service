package main

import (
	"context"
	pb "github.com/sigmawq/grpc-service/grpc"
	"github.com/sigmawq/grpc-service/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Sender struct{}

func (sender *Sender) Transmit(buffer []shared.DataEntry, host string) error {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to grpc server: %v", err)
		return err
	}
	defer conn.Close()

	grpcBuffer := make([]*pb.Data, len(buffer))
	for i, value := range buffer {
		grpcBuffer[i] = value.ToGrpcFormat()
	}

	context := context.Background()
	client := pb.NewGrpcClient(conn)
	_, err = client.SendBatch(context, &pb.Batch{Data: grpcBuffer})
	if err != nil {
		return err
	}

	return nil
}
