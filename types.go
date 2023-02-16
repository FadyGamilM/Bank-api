package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
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
