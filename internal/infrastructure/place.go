package infrastructure

import (
	"arch/internal/domain/entity"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PlaceGet struct {
	db *sqlx.DB
}

func NewPlace(db *sqlx.DB) *PlaceGet {
	return &PlaceGet{db: db}
}

func (r *PlaceGet) Get(params *entity.RequestPayload, centerLon, centerLat *float64) ([]entity.PlaceInfo, error) {
	var result []entity.PlaceInfo

	base := `
		SELECT title,address,description,lon,lat
		FROM place p
		LEFT JOIN place_schedule ps ON ps.place_id = p.id
		WHERE 1=1
`

	sql, args := builderSQL(params, base, centerLon, centerLat)
	sql += " GROUP BY p.id ORDER BY p.title LIMIT 50"

	if err := r.db.Select(&result, sql, args...); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return result, nil
}
