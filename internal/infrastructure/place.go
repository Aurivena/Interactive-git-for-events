package infrastructure

import (
	"arch/internal/domain/entity"
	"fmt"

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
	var output []entity.PlaceInfo

	base := `
		SELECT DISTINCT ON (p.id) id, title,kind,address,description,lon,lat,tags
		FROM place p
		LEFT JOIN place_schedule ps ON ps.place_id = p.id
		WHERE 1=1
`

	sql, args := builderSQL(params, base, centerLon, centerLat)
	sql += fmt.Sprintf(" ORDER BY p.id, RANDOM() LIMIT %d", params.Count)

	if err := r.db.Select(&output, sql, args...); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return output, nil
}

func (r *PlaceGet) ByID(id entity.UUID) (*entity.PlaceInfo, error) {
	var output entity.PlaceInfo

	if err := r.db.Get(&output, `SELECT id,title,kind,address,description,lon,lat,tags FROM place WHERE id = $1`, id); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &output, nil
}

func (r *PlaceGet) List() ([]entity.PlaceInfo, error) {
	var output []entity.PlaceInfo

	if err := r.db.Select(&output, `SELECT id,title,kind,address,description,lon,lat,tags FROM place`); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return output, nil
}

func (r *PlaceGet) ListByKind(kind entity.Kind) ([]entity.PlaceInfo, error) {
	var output []entity.PlaceInfo

	if err := r.db.Select(&output, `SELECT id,title,kind,address,description,lon,lat,tags FROM place WHERE kind = $1`, kind); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return output, nil
}
