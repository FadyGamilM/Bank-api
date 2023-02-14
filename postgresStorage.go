package main

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v4/stdlib"
	// _ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func PostgresStorageFactory() (*PostgresStorage, error) {
	DataSourceDetails := url.URL{
		Scheme: "postgres",
		Host:   "localhost:2345",
		User:   url.UserPassword("fady", "gobankingpassword"),
		Path:   "bankDB",
	}

	q := DataSourceDetails.Query()
	q.Add("sslmode", "disable")

	DataSourceDetails.RawQuery = q.Encode()

	// connectionString := "user=user dbname=bankDB password=gobankingpassword sslmode=disable"

	db, err := sql.Open("pgx", DataSourceDetails.String())

	defer func() {
		_ = db.Close()
		fmt.Println("Db Connection is closed")
	}()

	if err != nil {
		fmt.Println("error 1")
		return nil, err
	}

	// ping the db conn
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		fmt.Println("error 2")
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}
