package main

var database Database

func main() {
	_database, err := NewDatabase()
	if err != nil {
		return
	}
	database = _database

	_, err = NewServer("localhost:9000")
	if err != nil {
		return
	}
}
