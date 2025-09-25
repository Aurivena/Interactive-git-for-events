package ports

import (
	"arch/internal/domain/entity"
	"encoding/json"
)

type TourWrite interface {
	Write(dateFrom, dateTo, sessionID string, places entity.TourOutput) (*entity.UUID, error)
}

type TourReader interface {
	Reader(sessionID string) ([]entity.TourOutput, error)
	ReaderByID(id entity.UUID) (*entity.TourOutput, error)
}

type TourGenerate interface {
	GenerateTour(input entity.RouteParams, lon, lat float64) (json.RawMessage, error)
}
