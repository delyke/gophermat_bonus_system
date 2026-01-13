package jwtissuer

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Issuer struct {
	secret []byte
}

func New(secret string) *Issuer {
	return &Issuer{secret: []byte(secret)}
}

func (i *Issuer) IssueAccessToken(userUUID, login string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userUUID,
		"login": login,
		"iat":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(i.secret)
}
