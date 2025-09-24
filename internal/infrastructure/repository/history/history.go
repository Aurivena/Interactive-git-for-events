package history

import (
	"github.com/jmoiron/sqlx"
)

type History struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *History {
	return &History{db: db}
}
