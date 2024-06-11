package main

func main() {
	path := "data/data.json"

	rawBufferSize := 1 * 1024 * 1024
	maxObjects := 5
	parser, err := NewParserFromPath(path, rawBufferSize, maxObjects)
	if err != nil {
		return
	}

	sender, err := NewSender("localhost:9000")
	if err != nil {
		return
	}

	for parser.More() {
		err := parser.Parse()
		if err != nil {
			break
		}

		err = sender.SendBatch(parser.Buffer)
		if err != nil {
			return
		}
	}

}
