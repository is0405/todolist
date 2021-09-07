package model

type ToDO struct {
	ID          int    `db:"id" json:"id"`
	OK          bool   `db:"ok" json:"ok"`
	Title       string `db:"title" json:"title"`
	Memo        string `db:"memo" json:"memo"`
	AccountID   int    `db:"account_id" json:"account_id"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}
