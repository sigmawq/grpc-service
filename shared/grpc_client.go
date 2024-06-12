package shared

import (
	"context"
	"errors"
	pb "github.com/sigmawq/grpc-service/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client struct {
	client pb.GrpcClient
}

func NewClientFromHost(host string) (Client, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to grpc server: %v", err)
		return Client{}, err
	}

	client := pb.NewGrpcClient(conn)
	return NewClientFromGrpc(client)
}

func NewClientFromGrpc(client pb.GrpcClient) (Client, error) {
	sender := Client{
		client: client,
	}

	return sender, nil
}

func (client *Client) SendBatch(buffer []DataEntry) error {
	grpcBuffer := make([]*pb.Data, len(buffer))
	for i, value := range buffer {
		value := value.ToGrpcFormat()
		grpcBuffer[i] = &value
	}

	_, err := client.client.SendBatch(context.Background(), &pb.Batch{Data: grpcBuffer})
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) Retrieve(search string, from, size int32) ([]*pb.Data, error) {
	request := pb.RetrieveRequest{
		Search: search,
		From:   from,
		Size:   size,
	}
	resp, err := client.client.Retrieve(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New("failed to fetch data over GRPC")
	}

	return resp.Data, nil
}

func (client *Client) Aggregate() ([]*pb.AggregationCategory, error) {
	request := pb.AggregateRequest{}
	resp, err := client.client.Aggregate(context.Background(), &request)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, errors.New("failed to fetch data over GRPC")
	}

	return resp.Data, nil
}
