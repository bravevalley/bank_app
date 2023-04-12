// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"time"
)

type Account struct {
	AccNumber int64     `json:"acc_number"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID        int64     `json:"id"`
	AccNumber int64     `json:"acc_number"`
	Amount    int64     `json:"amount"`
	Date      time.Time `json:"date"`
}

type Transfer struct {
	ID     int64     `json:"id"`
	Amount int64     `json:"amount"`
	Debit  int64     `json:"debit"`
	Credit int64     `json:"credit"`
	Date   time.Time `json:"date"`
}

type User struct {
	Username            string    `json:"username"`
	HashedPassword      string    `json:"hashed_password"`
	FullName            string    `json:"full_name"`
	Email               string    `json:"email"`
	PasswordLastChanged time.Time `json:"password_last_changed"`
	CreatedAt           time.Time `json:"created_at"`
}
