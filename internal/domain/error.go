package domain

import "errors"

var (
	FileDuplicate = errors.New("file already exists")
)
