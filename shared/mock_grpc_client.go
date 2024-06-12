package shared

import (
	"context"
	pb "github.com/sigmawq/grpc-service/grpc"
	"google.golang.org/grpc"
)

type MockGrpcClient struct {
	SendBatchCount  int
	RetrieveCount   int
	AggregateCount  int
	LastInputLength int
}

func (client *MockGrpcClient) SendBatch(ctx context.Context, in *pb.Batch, opts ...grpc.CallOption) (*pb.BatchResponse, error) {
	client.SendBatchCount++
	client.LastInputLength = len(in.Data)
	return nil, nil
}

func (client *MockGrpcClient) Retrieve(ctx context.Context, in *pb.RetrieveRequest, opts ...grpc.CallOption) (*pb.RetrieveResponse, error) {
	client.RetrieveCount++
	client.LastInputLength = 0
	return nil, nil
}

func (client *MockGrpcClient) Aggregate(ctx context.Context, in *pb.AggregateRequest, opts ...grpc.CallOption) (*pb.AggregateResponse, error) {
	client.AggregateCount++
	client.LastInputLength = 0
	return nil, nil
}
