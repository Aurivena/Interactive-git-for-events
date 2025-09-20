package infrastructure

import (
	"arch/internal/domain/entity"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type History struct {
	db *sqlx.DB
}

func NewHistory(db *sqlx.DB) *History {
	return &History{db: db}
}

func (r *History) Save(aiMessage entity.ChatOutput, message string, sessionID string) error {
	_, err := r.db.Exec(`INSERT INTO  history (session,message,ai_message) VALUES ($1,$2,$3)`, sessionID, message, aiMessage)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r *History) ListBySessionID(query *entity.Query, session string) ([]entity.History, error) {
	var output []entity.History

	offset := (query.Page - 1) * query.Limit

	if err := r.db.Select(&output,
		`SELECT id, message, ai_message
					FROM history
					WHERE session = $1
					ORDER BY created_at
					LIMIT $2 OFFSET $3;`, session, query.Limit, offset); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return output, nil
}
