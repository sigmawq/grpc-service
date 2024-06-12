package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	pb "github.com/sigmawq/grpc-service/grpc"
	"github.com/sigmawq/grpc-service/shared"
	"log"
	"strings"
)

type Database struct {
	es *elasticsearch.Client
}

func NewDatabase(databaseHost string) (Database, error) {
	database := Database{}

	cfg := elasticsearch.Config{
		Addresses: []string{
			databaseHost,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to connect to elasticsearch cluser: %v", err)
	}

	database.es = es

	err = database.setupDatabaseIndex()
	if err != nil {
		return database, err
	}

	return database, nil
}

func (database *Database) setupDatabaseIndex() error {
	query := `
{
	"settings": {
        "analysis": {
            "analyzer": {
                "default": {
                    "type": "custom",
                    "tokenizer": "standard",
                    "filter": [
                        "lowercase",
                        "asciifolding"
                    ]
                }
            }
        }
  },
	"mappings": {
	    "properties": {
	      "subcategory": {
	        "type": "text",
	        "fielddata": true
	      }
	    }
	}
}`

	indexExistsReq := esapi.IndicesExistsRequest{
		Index: []string{"documents"},
	}

	resp, err := indexExistsReq.Do(context.Background(), database.es)
	if err != nil {
		log.Printf("Failed to check if index exists: %v", err)
		return err
	}

	if !resp.IsError() { // Index is already present
		log.Printf("Database index is already present")
		return nil
	}

	log.Printf("Database index does not exist, creating..")

	indexReq := esapi.IndicesCreateRequest{
		Index: "documents",
		Body:  strings.NewReader(query),
	}

	resp, err = indexReq.Do(context.Background(), database.es)
	if err != nil {
		log.Printf("Failed to define database mapping: %v", err)
		return err
	}

	if resp.IsError() {
		log.Printf("Failed to define database mapping: %v", resp)
		return errors.New("failed to define database mapping")
	}

	return nil
}

func (database *Database) UpdateBatch(batch []*pb.Data) error {
	_, err := database.es.Indices.Create("documents")
	if err != nil {
		log.Printf("Failed to create elasticsearch index: %v", err)
		return err
	}

	for _, dataGrpc := range batch {
		data := shared.NewDataEntryFromGrpcFormat(dataGrpc)

		serialized, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to serialize data: %v", err)
			return err
		}

		res, err := database.es.Index("documents", strings.NewReader(string(serialized)), database.es.Index.WithDocumentID(data.Id))
		if err != nil {
			log.Printf("Error while inserting documents: %v", err)
			return err
		}

		if res.IsError() {
			log.Printf("Error while inserting documents: %v", res)
			return errors.New("error while inserting documents")
		}
	}

	return nil
}

func (database *Database) Retrieve(search string, size, from int) ([]*shared.DataEntryDatabase, error) {
	// TODO: Erase diacritics

	query := fmt.Sprintf(`
{
    "query": {
        "multi_match": {
            "query": "%v",
            "fields": [
                "title_ro",
                "title_ru"
            ]
        }
    },
    "size": %v,
    "from": %v
}`, search, size, from)

	res, err := database.es.Search(
		database.es.Search.WithContext(context.Background()),
		database.es.Search.WithIndex("documents"),
		database.es.Search.WithBody(strings.NewReader(query)),
		database.es.Search.WithTrackTotalHits(true),
		database.es.Search.WithPretty(),
	)

	if err != nil {
		log.Printf("Search request failed: %v", res)
		return nil, err
	}

	if res.IsError() {
		log.Printf("Search query returned an error: %v", res)
		return nil, err
	}

	type Response struct {
		Hits struct {
			Hits []struct {
				Source shared.DataEntryDatabase `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Printf("Failed to decode search result: %v", err)
		return nil, err
	}

	result := make([]*shared.DataEntryDatabase, 0, len(response.Hits.Hits))
	for _, value := range response.Hits.Hits {
		result = append(result, &value.Source)
	}

	return result, nil
}

func (database *Database) Aggregate() ([]*shared.AggregationCategory, error) {
	query := `
{
   "size": 0, 
   "aggregations": {
      "aggregated": {
         "terms": {
            "field": "subcategory"
         }
      }
   }
}
`

	res, err := database.es.Search(
		database.es.Search.WithContext(context.Background()),
		database.es.Search.WithIndex("documents"),
		database.es.Search.WithBody(strings.NewReader(query)),
		database.es.Search.WithTrackTotalHits(true),
		database.es.Search.WithPretty(),
	)

	if err != nil {
		log.Printf("Search request failed: %v", res)
		return nil, nil
	}

	if res.IsError() {
		log.Printf("Search query failed: %v", res)
		return nil, nil
	}

	type Response struct {
		Aggregations struct {
			Aggregated struct {
				Buckets []shared.AggregationCategory `json:"buckets"`
			} `json:"aggregated"`
		} `json:"aggregations"`
	}

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Printf("Failed to decode search result: %v", err)
		return nil, nil
	}

	result := make([]*shared.AggregationCategory, 0, len(response.Aggregations.Aggregated.Buckets))
	for _, value := range response.Aggregations.Aggregated.Buckets {
		value := value
		result = append(result, &value)
	}

	return result, nil
}
