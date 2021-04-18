package api

import (
	"github.com/gin-gonic/gin"
)

// LoginIn representing the request body for POST /login
type LoginIn struct {
	Email    string `body:"email" validate:"required,min=1,email"`
	Password string `body:"password" validate:"required,min=1"`
}

type LoginOut struct {
	Token string `json:"token"`
}

//handleLogin is the controller of the route POST /login.
func (a *API) handleLogin(ctx *gin.Context, loginIn *LoginIn) (LoginOut, error) {
	jwt, err := a.loginService.Login(loginIn.Email, loginIn.Password)
	if err != nil {
		return LoginOut{}, err
	}

	return LoginOut{
		Token: jwt,
	}, nil
}
