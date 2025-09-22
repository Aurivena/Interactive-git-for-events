package ports

import "arch/internal/domain/entity"

type HistoryReader interface {
	ListBySessionID(query *entity.Query, session string) ([]entity.History, error)
}

type HistoryWriter interface {
	Write(aiMessage []entity.ChatOutput, message string, sessionID string) error
}
