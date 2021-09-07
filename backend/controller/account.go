package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"regexp"
	"errors"
	"fmt"
	// "strconv"
	"crypto/sha256"

	"github.com/is0405/model"
	"github.com/is0405/repository"
	"github.com/is0405/service"
	// "github.com/is0405/httputil"
	"github.com/jmoiron/sqlx"
	"encoding/hex"
)

type Account struct {
	db *sqlx.DB
}

func NewAccount(db *sqlx.DB) *Account {
	return &Account{db: db}
}

func (a *Account) Create(_ http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	
	account := &model.Account{}
	// request bodyの中にあるjso構造体にデコードする.
	err := json.NewDecoder(r.Body).Decode(&account);
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	fmt.Println(account.Mail)
	fmt.Println(account.Password)
	
	if account.Mail == "" {
		return http.StatusUnprocessableEntity, nil, errors.New("Not enough parameters")
	}

	if !MailCheck(account.Mail) {
		return http.StatusBadRequest, nil, errors.New("This email address is incorrect")
	}

	ok, err := repository.FindMail(a.db, account.Mail)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if ok != 0 {
		return http.StatusBadRequest, nil, errors.New("This email address is already registered")
	}

	account.Password = HashGenerateSha256(account.Password)
	
	Service := service.NewAccount(a.db)
	_, err = Service.Create(account)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	
	return http.StatusOK, nil, nil

}

func HashGenerateSha256(s string) string {
    b := []byte(s)
	sha256 := sha256.Sum256(b)

	return hex.EncodeToString(sha256[:])
}

func MailCheck(str string) bool {
	chars := []string{"@", ".", "\\_", "\\-"}
    r := strings.Join(chars, "")
	symbol := regexp.MustCompile("[^" + r + "A-Za-z0-9]+")
	if symbol.Match([]byte(str)) {
		//上記以外がある
		return false
	} else {
		symbol := regexp.MustCompile(`\s*@\s*`)
		symbol2 := regexp.MustCompile(`\s*\.\s*`)

		group := symbol.Split(str, -1)
		if len(group) != 2 {
			return false
		}

		group = symbol2.Split(str, -1)
		for i := 0; i < len(group); i++ {
			if group[i] == "" {
				return false
			} else if strings.HasSuffix(group[i], "@") {
				return false
			}
		}
	}
	return true
}
