package main

import "math/rand"

type Account struct {
	ID        int
	FirstName string
	LastName  string
	Number    int64
	Balance   int64
}

// initializer constructor
func AccountFactory(fname, lname string) *Account {
	return &Account{
		ID:        rand.Intn(500),
		FirstName: fname,
		LastName:  lname,
		Number:    int64(rand.Intn(100000000)),
	}
}
