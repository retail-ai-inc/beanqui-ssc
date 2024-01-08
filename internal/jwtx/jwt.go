package jwtx

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	UserName string
	jwt.RegisteredClaims
}

var signKey = []byte("sfds234@#$@4242")

func MakeRsaToken(claims Claim) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return str, nil
}
func ParseRsaToken(tokenStr string) (*Claim, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &Claim{}, func(token *jwt.Token) (i interface{}, err error) {
		return signKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*Claim); ok && token.Valid {
		return claim, nil
	}
	return nil, errors.New("invalid token")
}
