package main

import (
	"github.com/sigmawq/grpc-service/shared"
	"strings"
	"testing"
)

func TestParseAndSend(t *testing.T) {
	json := `
[
 {
  "_id": "38118545",
  "categories": {
   "subcategory": "1407"
  },
  "title": {
   "ro": "name ro 1",
   "ru": "name ru 1"
  },
  "type": "standard",
  "posted": 1486556302.101039
 },
 {
  "_id": "38784049",
  "categories": {
   "subcategory": "1404"
  },
  "title": {
   "ro": "name ro 2",
   "ru": "name ru 2"
  },
  "type": "standard",
  "posted": 1488274575.697526
 }
]
`

	bufferSize := 1 * 1024 * 1024
	maxObjects := 10000
	parser, err := NewParserFromReader(strings.NewReader(json), bufferSize, maxObjects)
	if err != nil {
		t.Errorf("Parser initialization failed")
	}

	grpcClient := shared.MockGrpcClient{}
	client, err := shared.NewClientFromGrpc(&grpcClient)
	if err != nil {
		t.Errorf("GRPC client initialization failed: %v", err)
	}

	ParseAndSend(&parser, &client)

	if grpcClient.SendBatchCount != 1 {
		t.Errorf("client.SendBatchCount is %v, but should be 1", grpcClient.SendBatchCount)
	}

	if grpcClient.LastInputLength != 2 {
		t.Errorf("client.LastInputLength is %v, but should be 3", grpcClient.LastInputLength)
	}

	return
}
