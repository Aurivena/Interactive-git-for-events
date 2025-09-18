package ports

import "arch/internal/domain/entity"

type PlaceGetter interface {
	Get(params *entity.RequestPayload, centerLon, centerLat *float64) ([]entity.PlaceInfo, error)
}
