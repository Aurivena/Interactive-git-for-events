package tour

import "github.com/jmoiron/sqlx"

type Tour struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Tour {
	return &Tour{db: db}
}
