package repository

import (
	"database/sql"
	"github.com/is0405/model"
	"github.com/jmoiron/sqlx"
)

func AddAccount(db *sqlx.DB, ma *model.Account) (sql.Result, error) {
	return db.Exec(`
INSERT INTO account (mail, password)
VALUES (?, ?)
`, ma.Mail, ma.Password)
}

func FindAccount(db *sqlx.DB, ma *model.Account) (int, error) {
	var a int
	if err := db.Get(&a, `
SELECT id
FROM account
WHERE mail = ? AND password = ?;
`, ma.Mail, ma.Password); err != nil {
		return 0, err
	}
	return a, nil
}

func FindMail(db *sqlx.DB, mail string) (int, error) {
	var a int
	if err := db.Get(&a, `
SELECT count(*)
FROM account
WHERE mail = ?;
`, mail); err != nil {
		return -1, err
	}
	return a, nil
}
