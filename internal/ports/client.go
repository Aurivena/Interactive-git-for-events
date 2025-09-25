package ports

import (
	"arch/internal/domain/entity"
	"encoding/json"
)

type ClientWrite interface {
	Write(sessionID string, survey entity.Survey) error
}

type ClientReade interface {
	Read(sessionID string) (json.RawMessage, error)
}
