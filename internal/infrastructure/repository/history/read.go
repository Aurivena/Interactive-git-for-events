package history

import (
	"arch/internal/domain/entity"

	"github.com/sirupsen/logrus"
)

func (r *History) ListBySessionID(query *entity.Query, session string) ([]entity.History, error) {
	var output []entity.History

	offset := (query.Page - 1) * query.Limit

	if err := r.db.Select(&output,
		`SELECT id, message, ai_message
				FROM history
				WHERE session = $1
				ORDER BY message, created_at DESC
					LIMIT $2 OFFSET $3;`, session, query.Limit, offset); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return output, nil
}
