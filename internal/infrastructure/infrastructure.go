package infrastructure

import (
	"arch/internal/application/ports"

	"github.com/jmoiron/sqlx"
)

type Infrastructure struct {
	PlaceGet ports.PlaceGetter
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func New(sources *Sources) *Infrastructure {
	return &Infrastructure{
		PlaceGet: NewPlace(sources.BusinessDB),
	}
}
