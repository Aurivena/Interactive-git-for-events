package client

import "github.com/jmoiron/sqlx"

type Client struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Client {
	return &Client{db: db}
}
