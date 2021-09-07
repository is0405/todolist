package controller

import (
	"encoding/json"
	"errors"
	//"fmt"
	"net/http"
	"time"

	"github.com/is0405/model"
	"github.com/is0405/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

type Login struct {
	db           *sqlx.DB
	jwtSecretKey []byte
}

func NewLogin(db *sqlx.DB, jwtSecretKey []byte) *Login {
	return &Login{db: db, jwtSecretKey: jwtSecretKey}
}

type Response struct {
	Token string `json:"token"`
}

func (a *Login) Login(_ http.ResponseWriter, r *http.Request) (int, interface{}, error) {

	account := &model.Account{}

	// request bodyの中にあるjsonをoffice構造体にデコードする.
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return http.StatusBadRequest, nil, err
	}

	if account.Mail == "" || account.Password == "" {
		return http.StatusUnprocessableEntity, nil, errors.New("required parameter is missing")
	}

	account.Password = HashGenerateSha256(account.Password)
	// 事業所の情報がとってこれるかでログイン判定.
	id, err := repository.FindAccount(a.db, account)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	// jwtの作成
	claims := model.Claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //トークン発行から24時間で使えなくなる.
			Issuer:    "todolist",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.jwtSecretKey)

	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	res := Response{
		Token:    signedToken,
	}

	return http.StatusOK, res, nil
}
