package service

import (
	"github.com/is0405/dbutil"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/is0405/model"
	"github.com/is0405/repository"
	//"fmt"
)

type ToDO struct {
	db *sqlx.DB
}

func NewToDO(db *sqlx.DB) *ToDO{
	return &ToDO{db}
}

func (a *ToDO) Create(mt *model.ToDO) (int64, error) {
	var createdId int64
	if err := dbutil.TXHandler(a.db, func(tx *sqlx.Tx) error {
		
		todo, err := repository.AddToDOList(a.db, mt)	
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

func (a *ToDO) Delete(todo_id int, account_id int) (int64, error) {
	var createdId int64
	if err := dbutil.TXHandler(a.db, func(tx *sqlx.Tx) error {
		
		model, err := repository.RemoveToDOList(a.db, todo_id, account_id)
		if err != nil {
			return err
		}
		
		if err := tx.Commit(); err != nil {
			return err
		}
		
		id, err := model.LastInsertId()
		
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

func (a *ToDO) Update(data *model.ToDO) (int64, error) {
	var createdId int64
	if err := dbutil.TXHandler(a.db, func(tx *sqlx.Tx) error {
		
		model, err := repository.UpdateToDOList(a.db, data)
		if err != nil {
			return err
		}
		
		if err := tx.Commit(); err != nil {
			return err
		}
		
		id, err := model.LastInsertId()
		
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
