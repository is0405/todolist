package repository

import (
	"database/sql"
	"github.com/is0405/model"
	"github.com/jmoiron/sqlx"
)

func AddToDOList(db *sqlx.DB, mt *model.ToDO) (sql.Result, error) {
	return db.Exec(`
INSERT INTO todo (title, memo, account_id, ok)
VALUES (?, ?, ?, ?)
`, mt.Title, mt.Memo, mt.AccountID, mt.OK)
}

func GetToDOList(db *sqlx.DB, account_id int) ([]model.ToDO, error) {
	a := make([]model.ToDO, 0)
	if err := db.Select(&a, `
SELECT id, ok, title, memo, account_id, created_at, updated_at
FROM todo
WHERE account_id = ?;
`, account_id); err != nil {
		return nil, err
	}
	return a, nil
}

func RemoveToDOList(db *sqlx.DB, todo_id int, account_id int) (sql.Result, error) {
	return db.Exec(`
DELETE FROM todo WHERE id = ? AND account_id = ?;
`, todo_id, account_id)
}

func UpdateToDOList(db *sqlx.DB, mt *model.ToDO) (sql.Result, error) {
	return db.Exec(`
UPDATE todo 
SET title = ?, memo = ?, ok = ? 
WHERE id = ? AND account_id = ?;
`, mt.Title, mt.Memo, mt.OK, mt.ID, mt.AccountID)
}
