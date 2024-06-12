package main

import (
	"github.com/sigmawq/grpc-service/shared"
	"log"
)

func ParseAndSend(parser *Parser, client *shared.Client) error {
	for parser.More() {
		err := parser.Parse()
		if err != nil {
			break
		}

		err = client.SendBatch(parser.Buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	path := "data/data.json"

	bufferSize := 1 * 1024 * 1024
	maxObjects := 10000
	parser, err := NewParserFromPath(path, bufferSize, maxObjects)
	if err != nil {
		return
	}

	client, err := shared.NewClientFromHost("localhost:9000")
	if err != nil {
		return
	}

	err = ParseAndSend(&parser, &client)
	if err != nil {
		log.Printf("Worker failed to parse and send data: %v", err)
		return
	}
}
