package main

import (
	"fmt"
	"log"
)

var sender Sender

func main() {
	_sender, err := NewSender("localhost:9000")
	if err != nil {
		log.Printf("Failed to initialize sender: %v", err)
		return
	}
	sender = _sender

	resp1, err := sender.Retrieve("Ёлка", 0, 10)
	resp2, err := sender.Aggregate()

	fmt.Println(resp1, resp2)

}
