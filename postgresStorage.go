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
		return nil, fmt.Errorf("error while opening a connection to the database")
	}

	// ping the db conn
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error while trying to ping the connected databse")
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
		return fmt.Errorf("error while executing the INSERT query")
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
		log.Println("error while quering all rows from database => ", err)
		return nil, fmt.Errorf("error while fetching all rows from database")
	}
	for rows.Next() {
		// define a destination so the `rows.Scan()` method can map the value from columns into the destination props
		// NOTE => i defined the account as an address because Scan maps the values into pointers
		account, _ := ScanRowIntoAccount(rows)
		allAccounts = append(allAccounts, account)
	}
	return allAccounts, nil
}

func (storage *PostgresStorage) GetById(accountID int) (*Account, error) {
	query := `SELECT * FROM ACCOUNTS AS ACC WHERE ACC.id = $1`

	rows, err := storage.db.Query(query, accountID)
	if err != nil {
		log.Println("Error while fetching the account => ", err)
		return nil, err
	}

	for rows.Next() {
		return ScanRowIntoAccount(rows)
	}

	// return not found
	return nil, fmt.Errorf("Account with id = %d is not found", accountID)
}

// ! a helper method used to scan one row [after calling .Next() from the caller function] and convert it into Account entity
func ScanRowIntoAccount(rows *sql.Rows) (*Account, error) {
	// define new entity to be a placeholder
	account := new(Account)

	// scan the current row into this entity
	err := rows.Scan(&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	if err != nil {
		log.Println("Error while scanning the row => ", err)
		return nil, fmt.Errorf("error while scanning the fetched row")
	}
	return account, nil
}

func (storage *PostgresStorage) Update(account *Account) error {
	return nil
}

func (storage *PostgresStorage) DeleteById(accountID int) error {
	sqlResult, err := storage.db.Exec(`DELETE FROM ACCOUNTS WHERE id = $1`, accountID)
	if err != nil {
		return fmt.Errorf("error while deleting the account with id = %d", accountID)
	}
	log.Printf("Result of deleting account with id = %d \n", accountID)
	log.Println(sqlResult.RowsAffected())
	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("server Error while deleting a row with id = %d", accountID)
	}
	log.Printf("Number of affected rows is %d", rowsAffected)
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
