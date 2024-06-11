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

func NewDatabase() (Database, error) {
	database := Database{}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9500",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to connect to elasticsearch cluser: %v", err)
	}

	database.es = es

	err = database.allowAggregationOnSubcategoryField()
	if err != nil {
		return database, err
	}

	//ind, err := es.Indices.Create("documents")
	//if err != nil {
	//	return database, err
	//}
	//log.Println(ind)
	//
	//res, err := es.Index("documents", strings.NewReader(`{"title" : "Test"}`), es.Index.WithDocumentID("1"))
	//log.Println(res)
	//
	//res, err = es.Get("documents", "1")
	//
	//res.
	//
	//buffer, err := io.ReadAll(res.Body)
	//log.Println(buffer)

	return database, nil
}

func (database *Database) allowAggregationOnSubcategoryField() error {
	query := `
{
"mappings": {
	    "properties": {
	      "subcategory": {
	        "type": "text",
	        "fielddata": true
	      }
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
		return errors.New("Failed to define database mapping")
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
			log.Printf("")
			return err
		}
		log.Println(res)
	}

	return nil
}

func (database *Database) Retrieve(search string, size, from int) ([]interface{}, error) {
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
		log.Printf("Search query failed: %v", res)
		return nil, err
	}

	var bodyResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&bodyResult)
	if err != nil {
		log.Printf("Failed to decode search result: %v", err)
		return nil, err
	}

	values, _ := bodyResult["hits"].(map[string]interface{})["hits"].([]interface{})

	return values, nil
}

func (database *Database) Aggregate() ([]interface{}, error) {
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

	var bodyResult map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&bodyResult)
	if err != nil {
		log.Printf("Failed to decode search result: %v", err)
		return nil, nil
	}

	values, _ := bodyResult["aggregations"].(map[string]interface{})["aggregated"].(map[string]interface{})["buckets"].([]interface{})

	return values, nil
}
