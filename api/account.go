package api

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"time"
)

// AccountCreateIn represents the JSON object received for an account creation.
type AccountCreateIn struct {
	// Email is the account email
	Email string `body:"email" validate:"required,min=1,email"`
	// Password is the account password
	Password string `body:"password" validate:"required,min=8"`
}

// AccountUpdateIn represents the data received
// to update the account corresponding to the given ID
type AccountUpdateIn struct {
	// ID is the account ID
	ID int `path:"id" validate:"required,min=1"`
	// Email is the new account email
	Email string `body:"email" validate:"omitempty,min=1,email"`
	// Password is the new account password
	Password string `body:"password" validate:"omitempty,min=8"`
}

// AccountDeleteIn contains the required field used to delete an existing account.
type AccountDeleteIn struct {
	// ID is the account ID to delete
	ID int `path:"id" validate:"required,min=1"`
}

// AccountOut represents the JSON object sent after an account creation.
type AccountOut struct {
	// ID is the account ID
	ID int `json:"id"`
	// Email is the account email
	Email string `json:"email"`
	// CreatedAt is the timestamp at which the account has been created
	CreatedAt time.Time `json:"createdAt"`
}

//createAccount is the controller used to create an accout.
func (a *API) createAccount(ctx *gin.Context, in *AccountCreateIn) (AccountOut, error) {
	createdAccount, err := a.accountService.CreateAccount(in.Email, in.Password)
	if err != nil {
		return AccountOut{}, err
	}

	return AccountOut{
		ID:        createdAccount.ID,
		Email:     createdAccount.Email,
		CreatedAt: createdAccount.CreatedAt,
	}, nil
}

//updateAccount is the controller used to update an existing account.
func (a *API) updateAccount(ctx *gin.Context, in *AccountUpdateIn) (AccountOut, error) {
	if len(in.Email) == 0 && len(in.Password) == 0 {
		return AccountOut{}, errors.NotValidf("fields")
	}

	updatedAccount, err := a.accountService.UpdateAccount(in.ID, in.Email, in.Password)
	if err != nil {
		return AccountOut{}, err
	}

	return AccountOut{
		ID:        updatedAccount.ID,
		Email:     updatedAccount.Email,
		CreatedAt: updatedAccount.CreatedAt,
	}, nil
}

func (a *API) deleteAccount(ctx *gin.Context, in *AccountDeleteIn) error {
	return a.accountService.DeleteAccount(in.ID)
}
