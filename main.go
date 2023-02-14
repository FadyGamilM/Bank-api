package main

import (
	"fmt"
	"log"
)

func main() {
	// connect to the database
	storage, err := PostgresStorageFactory()
	if err != nil {
		log.Fatalf("Can't connect to the persistence provider => %s", err)
	}

	fmt.Printf("%v\n", storage)

	// run the server
	server := NewApiServer("localhost:5000", storage)
	server.Run()
}
