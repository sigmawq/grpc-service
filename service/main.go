package main

import (
	"log"
	"os"
)

var db Database

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	elasticHost := os.Getenv("ELASTIC_HOST")
	if elasticHost == "" {
		elasticHost = "http://localhost:9200"
	}

	log.Printf("PORT=%v, ELASTIC_HOST=%v", port, elasticHost)

	_db, err := NewDatabase(elasticHost)
	if err != nil {
		os.Exit(1)
	}
	db = _db

	_, err = NewServer(":" + port)
	if err != nil {
		os.Exit(1)
	}
}
