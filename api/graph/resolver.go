package graph

import (
	"fmt"
	"github.com/sigmawq/grpc-service/shared"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var client shared.Client

func InitializeGraphQLClient(serviceHost string) error {
	_client, err := shared.NewClientFromHost(serviceHost)
	if err != nil {
		fmt.Printf("Failed to initialize sender: %v", err)
		return err
	}
	client = _client

	return err
}

type Resolver struct{}
