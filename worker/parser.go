package main

import (
	"bufio"
	"encoding/json"
	pb "github.com/sigmawq/grpc-service/grpc"
	"log"
	"os"
)

type Parser struct {
	Buffer []DataEntry
}

type DataEntry struct {
	Id         string `json:"_id"`
	Categories struct {
		Subcategory string `json:"subcategory"`
	} `json:"categories"`
	Title struct {
		Ro string `json:"ro"`
		Ru string `json:"ru"`
	} `json:"title"`
	Type   string  `json:"type"`
	Posted float64 `json:"posted"`
}

func (de *DataEntry) ToGrpcFormat() *pb.Data {
	return &pb.Data{
		Id:          de.Id,
		Subcategory: de.Categories.Subcategory,
		TitleRo:     de.Title.Ro,
		TitleRu:     de.Title.Ru,
		Type:        de.Type,
		Posted:      de.Posted,
	}
}

func (parser *Parser) ParseFromPath(path string, rawBufferSize, maxObjects int) error {
	data, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to read data.json")
	}

	reader := bufio.NewReaderSize(data, rawBufferSize)

	parser.Buffer = make([]DataEntry, 0)

	decoder := json.NewDecoder(reader)

	token, err := decoder.Token()
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	delim, _ := token.(json.Delim)
	if delim != '[' {
		log.Printf("Failed to decode, Expected '[', but got '%v'", token)
		return err
	}

	var object DataEntry
	for decoder.More() {
		if err != nil {
			log.Printf("%v", err)
			break
		}

		err = decoder.Decode(&object)
		if err != nil {
			log.Printf("Error while decoding, %v", err)
			return err
		}

		parser.Buffer = append(parser.Buffer, object)
		if len(parser.Buffer) > maxObjects {
			return nil
		}
	}

	return nil
}
