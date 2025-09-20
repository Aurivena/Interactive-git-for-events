package entity

import "encoding/json"

type Query struct {
	Page  int `json:"page" minimum:"1"`
	Limit int `json:"limit" minimum:"10"`
}

type History struct {
	ID        string          `json:"id" db:"id"`
	Message   string          `json:"message" db:"message"`
	AiMessage json.RawMessage `json:"ai_message" db:"ai_message"`
}
