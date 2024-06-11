package graph

import "fmt"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var sender Sender

func InitializeGraphQLSender() error {
	_sender, err := NewSender("localhost:9000")
	if err != nil {
		fmt.Printf("Failed to initialize sender: %v", err)
		return err
	}
	sender = _sender

	return err
}

type Resolver struct{}
