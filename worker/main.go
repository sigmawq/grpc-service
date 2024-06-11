package main

func main() {
	path := "data/data.json"

	rawBufferSize := 1 * 1024 * 1024
	maxObjects := 5
	parser, err := NewParserFromPath(path, rawBufferSize, maxObjects)
	if err != nil {
		return
	}

	sender := Sender{}
	for parser.More() {
		err := parser.Parse()
		if err != nil {
			break
		}

		err = sender.Transmit(parser.Buffer, "localhost:9000")
		if err != nil {
			return
		}
	}

}
