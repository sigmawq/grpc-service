package main

import (
	"bufio"
	"encoding/json"
	"github.com/sigmawq/grpc-service/shared"
	"io"
	"log"
	"os"
)

type Parser struct {
	file       io.Reader
	decoder    *json.Decoder
	maxObjects int
	hasMore    bool
	Buffer     []shared.DataEntry
}

func NewParserFromPath(path string, bufferSize, maxObjects int) (Parser, error) {
	data, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to read data: %v", err)
		return Parser{}, err
	}

	return NewParserFromReader(data, bufferSize, maxObjects)
}

func NewParserFromReader(reader io.Reader, bufferSize, maxObjects int) (Parser, error) {
	parser := Parser{}

	parser.file = bufio.NewReaderSize(reader, bufferSize)
	parser.decoder = json.NewDecoder(parser.file)
	parser.maxObjects = maxObjects
	parser.hasMore = true

	err := parser.consumeFirst()
	if err != nil {
		log.Printf("Failed to consume initial array bracket, data source is likely in an invalid format: %v", err)
		return parser, err
	}

	return parser, nil
}

func (parser *Parser) Parse() error {
	parser.Buffer = make([]shared.DataEntry, 0)

	var object shared.DataEntry
	for parser.decoder.More() {
		err := parser.decoder.Decode(&object)
		if err != nil {
			log.Printf("Error while decoding, %v", err)
			return err
		}

		parser.Buffer = append(parser.Buffer, object)
		if len(parser.Buffer) >= parser.maxObjects {
			return nil
		}
	}

	parser.hasMore = false
	return nil
}

func (parser *Parser) consumeFirst() error {
	token, err := parser.decoder.Token()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	delim, _ := token.(json.Delim)
	if delim != '[' {
		log.Printf("Failed to decode, Expected '[', but got '%v'", token)
		return err
	}

	return nil
}

func (parser *Parser) More() bool {
	return parser.hasMore
}
