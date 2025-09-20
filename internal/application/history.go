package application

import "arch/internal/domain/entity"

func (a *Application) ListHistory(query *entity.Query, sessionID string) ([]entity.History, error) {
	output, err := a.post.History.ListBySessionID(query, sessionID)
	if err != nil {
		return nil, err
	}
	return output, nil
}
