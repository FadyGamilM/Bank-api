package main

// abstract our persistence layer so we can later use any database provider [sql or nosql]
type Storage interface {
	CreateAccount(*Account) error

	DeleteAccount(int) error

	UpdateAccount(*Account) error

	GetAccountById(int) (*Account, error)
}
