package ports

import (
	"arch/internal/domain/entity"

	"github.com/google/uuid"
)

type PlaceReader interface {
	Get(params *entity.RequestPayload, centerLon, centerLat *float64) ([]entity.PlaceInfo, error)
	ByID(id entity.UUID) (*entity.PlaceInfo, error)
	ListByKind(kind entity.Kind) ([]entity.PlaceInfo, error)
	List() ([]entity.PlaceInfo, error)
}

type PlaceWriter interface {
	Write(id uuid.UUID, sql string) error
}

type PlaceBinding interface {
	Bind(placeID uuid.UUID, imageID string) error
}
