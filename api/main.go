package main

import (
	"log"
	"os"
)

func main() {
	data, err := os.Open("../data/data.json")
	if err != nil {
		log.Fatal("Failed to read data.json")
	}

	// jsonDataReader := strings.NewReader(jsonData)
	// decoder := json.NewDecoder(jsonDataReader)
	// var profile map[string]interface{}
	// err := decoder.Decode(&profile)
	// if err != nil {
	//     panic(err)
	// }
}
