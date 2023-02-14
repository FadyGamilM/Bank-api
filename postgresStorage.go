package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func PostgresStorageFactory() (*PostgresStorage, error) {
	connectionString := "user=postgres dbname=bankDB password=gobankingpassword sslmode=verify-full"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// ping the db conn
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}
