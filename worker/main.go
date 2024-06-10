package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Printf("Expected path to data file")
		return
	}
	path := os.Args[1]

	rawBufferSize := 1 * 1024 * 1024
	maxObjects := 10000
	parser := Parser{}
	err := parser.ParseFromPath(path, rawBufferSize, maxObjects)
	if err != nil {
		return
	}

	sender := Sender{}
	err = sender.Transmit(parser.Buffer, "localhost:9000")
	if err != nil {
		return
	}
}
