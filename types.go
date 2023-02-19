package main

import (
	"math/rand"
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"fname"`
	LastName  string    `json:"lname"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// initializer constructor
func AccountFactory(fname, lname string) *Account {
	rand.Seed(time.Now().Unix())
	return &Account{
		FirstName: fname,
		LastName:  lname,
		Number:    int64(rand.Intn(15000)),
		CreatedAt: time.Now().UTC(),
	}
}

// create account dto
type CreateAccountDTO struct {
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
}

// transfer to account request
type TransferToAccountRequest struct {
	ToAccountID int     `json:"to_account_id"`
	Amount      float64 `json:"amount"`
}
