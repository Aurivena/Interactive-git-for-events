package ports

import "arch/internal/domain/entity"

type ClientWrite interface {
	Write(sessionID string, survey entity.Survey) error
}
