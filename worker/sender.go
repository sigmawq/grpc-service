package main

import (
	"context"
	pb "github.com/sigmawq/grpc-service/grpc"
	"github.com/sigmawq/grpc-service/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Sender struct {
	client pb.GrpcClient
}

func NewSender(host string) (Sender, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to grpc server: %v", err)
		return Sender{}, err
	}

	sender := Sender{
		client: pb.NewGrpcClient(conn),
	}

	return sender, nil

}

func (sender *Sender) SendBatch(buffer []shared.DataEntry) error {
	grpcBuffer := make([]*pb.Data, len(buffer))
	for i, value := range buffer {
		grpcBuffer[i] = value.ToGrpcFormat()
	}

	context := context.Background()
	_, err := sender.client.SendBatch(context, &pb.Batch{Data: grpcBuffer})
	if err != nil {
		return err
	}

	return nil
}
