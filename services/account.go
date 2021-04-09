package services

import (
	"regexp"

	"github.com/fberrez/simple-uptime-backend/models"
	"github.com/juju/errors"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	repository models.AccountRepository
}

// NewAccountService intializes a new instance of this service and returns it.
func NewAccountService(repository models.AccountRepository) *Account {
	return &Account{
		repository: repository,
	}
}

func (a *Account) CreateAccount(email, password string) (*models.Account, error) {
	if !isEmailValid(email) {
		return nil, errors.NotValidf("email")
	}

	if !isPasswordValid(password) {
		return nil, errors.NotValidf("password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Annotatef(err, "generate from password")
	}

	createdAccount, err := a.repository.CreateAccount(email, string(hashedPassword))
	if err != nil {
		return nil, errors.Annotatef(err, "create an account")
	}

	return createdAccount, nil
}

func (a *Account) UpdateAccount(id int, email, password string) (*models.Account, error) {
	var updatedEmail string
	var updatedHashedPassword []byte
	var err error
	if len(email) > 0 {
		if !isEmailValid(email) {
			return nil, errors.NotValidf("email")
		}

		updatedEmail = email
	}

	if len(password) > 0 {
		if !isPasswordValid(password) {
			return nil, errors.NotValidf("password")
		}

		updatedHashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.Annotatef(err, "generate from password")
		}
	}

	updatedAccount, err := a.repository.UpdateAccount(id, updatedEmail, string(updatedHashedPassword))
	if err != nil {
		return nil, errors.Annotatef(err, "create an account")
	}

	return updatedAccount, nil
}

func (a *Account) DeleteAccount(id int) error {
	return a.repository.DeleteAccount(id)
}

// isEmailValid returns true if the given email is valid.
func isEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

// isPasswordValid returns true if the given password is valid.
func isPasswordValid(password string) bool {
	return len(password) > 8
}
