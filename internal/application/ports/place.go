package ports

import (
	"arch/internal/domain/entity"
)

type PlaceGetter interface {
	Get(params *entity.RequestPayload, centerLon, centerLat *float64) ([]entity.PlaceInfo, error)
	ByID(id entity.UUID) (*entity.PlaceInfo, error)
	ListByKind(kind entity.Kind) ([]entity.PlaceInfo, error)
	List() ([]entity.PlaceInfo, error)
}
