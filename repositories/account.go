package repositories

import (
	"github.com/fberrez/simple-uptime-backend/backend"
	"github.com/fberrez/simple-uptime-backend/models"
	"github.com/juju/errors"
)

// AccountRepository represents the account repository which contains backend instance.
type AccountRepository struct {
	backend backend.Backend
}

// NewAccountRepository initializes a new AccountRepo with the given backend instance.
func NewAccountRepository(backend backend.Backend) *AccountRepository {
	return &AccountRepository{
		backend: backend,
	}
}

// CreateAccount creates a new account with the given data and saves it in the backend instance.
func (a *AccountRepository) CreateAccount(email, password string) (*models.Account, error) {
	err := a.backend.ExecTransaction("INSERT INTO account(email, password) VALUES ($1, $2)", email, password)
	if err != nil {
		return nil, errors.Annotatef(err, "insert new account in backend")
	}

	return a.GetAccountByEmail(email)
}

// UpdateAccount updates the existing account corresponding to the given id
// with the given data.
func (a *AccountRepository) UpdateAccount(id int, email, password string) (*models.Account, error) {
	savedAccount, err := a.GetAccountByID(id)
	if err != nil {
		return nil, errors.Annotatef(err, "get account by id")
	}

	var emailToSave string
	if emailToSave = savedAccount.Email; len(email) > 0 && emailToSave != email {
		emailToSave = email
	}

	var passwordToSave string
	if passwordToSave = savedAccount.Password; len(password) > 0 && passwordToSave != password {
		passwordToSave = password
	}

	err = a.backend.ExecTransaction("UPDATE account SET email = $1, password = $2 WHERE id = $3", emailToSave, passwordToSave, id)
	if err != nil {
		return nil, errors.Annotatef(err, "update account")
	}

	return a.GetAccountByID(id)
}

// DeleteAccount builds and executes a query to delete the account corresponding to the given id.
func (a *AccountRepository) DeleteAccount(id int) error {
	return a.backend.ExecTransaction("DELETE FROM account where id = $1", id)
}

// GetAccountByEmail builds and executes a query to get the account corresponding to the given email.
func (a *AccountRepository) GetAccountByEmail(email string) (*models.Account, error) {
	account := models.Account{}
	err := a.backend.QueryRow("SELECT id, email, password, created_at FROM account WHERE email = $1", email).Scan(&account.ID, &account.Email, &account.Password, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// getAccountByID builds and executes a query to get the account corresponding to the given ID.
func (a *AccountRepository) GetAccountByID(id int) (*models.Account, error) {
	account := models.Account{}
	err := a.backend.QueryRow("SELECT id, email, password, created_at FROM account WHERE id = $1", id).Scan(&account.ID, &account.Email, &account.Password, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
