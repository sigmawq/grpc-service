package main

import (
	pb "github.com/sigmawq/grpc-service/grpc"
	"strconv"
)

func main() {
	database, err := NewDatabase()
	if err != nil {
		return
	}

	data := make([]*pb.Data, 5)
	for i, _ := range data {
		data[i] = new(pb.Data)
		v := data[i]

		id := strconv.FormatInt(int64(i)+100, 10)
		v.Id = id
		v.TitleRu = "Ёлка"
		v.TitleRo = "Apa"
		if i <= 2 {
			v.Subcategory = "subcategory1"
		} else {
			v.Subcategory = "subcategory2"
		}
	}
	err = database.UpdateBatch(data)
	if err != nil {
		return
	}

	database.Retrieve("Apa", 10, 0)

	database.Aggregate()

	_, err = NewServer("localhost:9000")
	if err != nil {
		return
	}

}
