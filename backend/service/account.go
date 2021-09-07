package service

import (
	"github.com/is0405/dbutil"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/is0405/model"
	"github.com/is0405/repository"
	//"fmt"
)

type Account struct {
	db *sqlx.DB
}

func NewAccount(db *sqlx.DB) *Account{
	return &Account{db}
}

func (a *Account) Create(ma *model.Account) (int64, error) {
	var createdId int64
	if err := dbutil.TXHandler(a.db, func(tx *sqlx.Tx) error {
		
		todo, err := repository.AddAccount(a.db, ma)	
		if err != nil {
			return err
		}
		
		if err := tx.Commit(); err != nil {
			return err
		}
		
		id, err := todo.LastInsertId()
		
		if err != nil {
			return err
		}
		
		createdId = id
		return err
		
	}); err != nil {
		return 0, errors.Wrap(err, "failed recipe insert transaction")
	}
	return createdId, nil
}
