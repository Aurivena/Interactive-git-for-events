package ports

import "arch/internal/domain/entity"

type History interface {
	Save(aiMessage []entity.ChatOutput, message string, sessionID string) error
	ListBySessionID(query *entity.Query, session string) ([]entity.History, error)
}
