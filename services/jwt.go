package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/juju/errors"
	"github.com/ovh/configstore"
)

// Claim contains data that will be included in the jwt.
type Claim struct {
	// Email the account email related to the jwt
	Email string `json:"email"`
	jwt.StandardClaims
}

// jwtTTLMin determines the number of minutes in which the token is valid.
var jwtTTLMin = time.Minute * 5

// secretKey is the key corresponding to the secret value in config file, used to sign jwt
var secretKey = "jwtSecret"

// New generates a new signed jwt
func NewJWT(email string) (string, error) {
	claim := Claim{
		email,
		jwt.StandardClaims{
			Audience:  "authenticatedUser",
			ExpiresAt: time.Now().Add(jwtTTLMin).Unix(),
			Issuer:    "simpleUptime",
			Subject:   "useOfSecuredRoutes",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return SignJWT(token)
}

// SignJWT signs the given token.
func SignJWT(token *jwt.Token) (string, error) {
	secret, err := configstore.GetItemValue(secretKey)
	if err != nil {
		return "", errors.Annotatef(err, "signin jwt")
	}
	jwtToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.Annotatef(err, "signin jwt")
	}
	return jwtToken, nil
}
