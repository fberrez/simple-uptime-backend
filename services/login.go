package services

import (
	"github.com/fberrez/simple-uptime-backend/models"
	"github.com/juju/errors"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	accountRepository models.AccountRepository
}

// NewAccountService intializes a new instance of this service and returns it.
func NewLoginService(repository models.AccountRepository) *Login {
	return &Login{
		accountRepository: repository,
	}
}

// Login performs an authentication with the given data.
func (l *Login) Login(email, password string) (string, error) {
	foundAccount, err := l.accountRepository.GetAccountByEmail(email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(foundAccount.Password), []byte(password)); err != nil {
		return "", errors.NotValidf("password")
	}

	return NewJWT(email)
}
