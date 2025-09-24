package place

import "github.com/jmoiron/sqlx"

type Place struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Place {
	return &Place{db: db}
}
