package main

import (
	"database/sql"
	"fmt"
	"log"
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

	// connectionString := "user=fady dbname=bankDB password=gobankingpassword sslmode=disable"

	db, err := sql.Open("pgx", DataSourceDetails.String())
	// defer func() {
	// 	_ = db.Close()
	// 	fmt.Println("Db Connection is closed")
	// }()
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

// TODO => implement the storage interface
func (storage *PostgresStorage) Create(account *Account) error {
	query := `INSERT INTO ACCOUNTS 
		(fname, lname, number, balance, created_at)
		VALUES 
		($1, $2, $3, $4, $5)`
	queryResp, err := storage.db.Exec(query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
	)
	if err != nil {
		log.Println("Error while executing the INSERT query => ", err)
		return err
	}

	log.Println(queryResp)
	return nil
}

func (storage *PostgresStorage) GetAll() ([]*Account, error) {
	// the list we will return to the handler
	allAccounts := []*Account{}

	// execute the query and process the resulted data
	query := `SELECT * FROM ACCOUNTS`
	rows, err := storage.db.Query(query)
	if err != nil {
		log.Println("Error while quering all rows from database => ", err)
		return nil, err
	}
	for rows.Next() {
		// define a destination so the `rows.Scan()` method can map the value from columns into the destination props
		// NOTE => i defined the account as an address because Scan maps the values into pointers
		account := new(Account)
		err := rows.Scan(&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt)
		if err != nil {
			log.Println("Error while scanning the row => ", err)
			return nil, err
		}
		allAccounts = append(allAccounts, account)
	}
	return allAccounts, nil
}

func (storage *PostgresStorage) GetById(accountID int) (*Account, error) {
	return nil, nil
}

func (storage *PostgresStorage) Update(account *Account) error {
	return nil
}

func (storage *PostgresStorage) DeleteById(accountID int) error {
	return nil
}

func (storage *PostgresStorage) CreateAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS ACCOUNTS (
			id SERIAL PRIMARY KEY, 
			fname VARCHAR(50) NOT NULL,
			lname VARCHAR(50) NOT NULL,
			number serial,
			balance serial,
			created_at timestamp
		);`

	sqlResult, err := storage.db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println(sqlResult)
	return nil
}

func (storage *PostgresStorage) Init() error {
	return storage.CreateAccountTable()
}
