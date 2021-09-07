package model
type Account struct {
	ID          int    `db:"id" json:"id"`
	Mail        string `db:"mail" json:"mail"`
	Password    string `db:"password" json:"password"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}
