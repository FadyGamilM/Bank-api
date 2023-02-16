package main

// abstract our persistence layer so we can later use any database provider [sql or nosql]
type Storage interface {
	Create(*Account) error

	DeleteById(int) error

	Update(*Account) error

	GetById(int) (*Account, error)

	GetAll() ([]*Account, error)
}
