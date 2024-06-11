package main

import (
	"context"
	"errors"
	pb "github.com/sigmawq/grpc-service/grpc"
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

func (sender *Sender) Retrieve(search string, from, size int32) (string, error) {
	request := pb.RetrieveRequest{
		Search: search,
		From:   from,
		Size:   size,
	}
	resp, err := sender.client.Retrieve(context.Background(), &request)
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", errors.New("failed to fetch data over GRPC")
	}

	return resp.Data, nil
}

func (sender *Sender) Aggregate() (string, error) {
	request := pb.AggregateRequest{}
	resp, err := sender.client.Aggregate(context.Background(), &request)
	if err != nil {
		return "", err
	}

	if !resp.Success {
		return "", errors.New("failed to fetch data over GRPC")
	}

	return resp.Data, nil
}
