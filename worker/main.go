package main

import (
	"github.com/sigmawq/grpc-service/shared"
	"log"
	"os"
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
	dataSource := os.Getenv("DATA_SOURCE")
	if dataSource == "" {
		dataSource = "data/data.json"
	}

	serviceHost := os.Getenv("SERVICE_HOST")
	if serviceHost == "" {
		serviceHost = "localhost:9000"
	}

	log.Printf("DATA_SOURCE=%v, SERVICE_HOST=%v", dataSource, serviceHost)

	bufferSize := 1 * 1024 * 1024
	maxObjects := 10000
	parser, err := NewParserFromPath(dataSource, bufferSize, maxObjects)
	if err != nil {
		os.Exit(1)
	}

	client, err := shared.NewClientFromHost(serviceHost)
	if err != nil {
		os.Exit(1)
	}

	err = ParseAndSend(&parser, &client)
	if err != nil {
		log.Printf("Worker failed to parse and send data: %v", err)
		os.Exit(1)
	}

	log.Print("done!")
}
