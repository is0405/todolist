package model

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID                  int    `json:"account_id"`
	jwt.StandardClaims 
}
