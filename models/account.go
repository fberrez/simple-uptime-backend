package models

import (
	"time"
)

// Account represents an account.
type Account struct {
	// ID is the account ID
	ID int
	// Email is the account email
	Email string
	// Password is the account password
	Password string
	// CreatedAt is the time at which the user has been created.
	CreatedAt time.Time
	// LastConnection is the time at which the user last logged on.
	LastConnection time.Time
}

// AccountRepository is the interface containing function to interact
// with accounts saved in a backend instance
type AccountRepository interface {
	CreateAccount(email, password string) (*Account, error)
	UpdateAccount(id int, email, password string) (*Account, error)
	DeleteAccount(id int) error
}
