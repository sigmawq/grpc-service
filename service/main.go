package main

var db Database

func main() {
	_db, err := NewDatabase("http://localhost:9500")
	if err != nil {
		return
	}
	db = _db

	_, err = NewServer("localhost:9000")
	if err != nil {
		return
	}
}
