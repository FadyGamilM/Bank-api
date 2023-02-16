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

	// create a table for ACCOUNT entity
	if err := storage.Init(); err != nil {
		log.Fatal("Erro while creating a table => ", err)
	}

	// run the server
	server := NewApiServer("localhost:5000", storage)
	server.Run()
}
