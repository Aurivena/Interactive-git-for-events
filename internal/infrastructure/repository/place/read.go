package place

import (
	"arch/internal/domain/entity"
	"arch/pkg/builder"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (r *Place) Get(params *entity.RequestPayload, centerLon, centerLat *float64) ([]entity.PlaceInfo, error) {
	var output []entity.PlaceInfo

	base := `
		SELECT DISTINCT ON (p.id) id, title,kind,address,description,lon,lat,tags
		FROM place p
		WHERE 1=1
	`

	sql, args := builder.BuildSql(params, base, centerLon, centerLat)
	sql += fmt.Sprintf(" ORDER BY p.id, RANDOM() LIMIT %d", params.Count)

	if err := r.db.Select(&output, sql, args...); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return output, nil
}

func (r *Place) ByID(id entity.UUID) (*entity.PlaceInfo, error) {
	var output entity.PlaceInfo

	if err := r.db.Get(&output, `SELECT id,title,kind,address,description,lon,lat,tags FROM place WHERE id = $1`, id); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &output, nil
}

func (r *Place) List() ([]entity.PlaceInfo, error) {
	var output []entity.PlaceInfo

	if err := r.db.Select(&output, `SELECT id,title,kind,address,description,lon,lat,tags FROM place`); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return output, nil
}

func (r *Place) ListByKind(kind entity.Kind) ([]entity.PlaceInfo, error) {
	var output []entity.PlaceInfo

	if err := r.db.Select(&output, `SELECT id,title,kind,address,description,lon,lat,tags FROM place WHERE kind = $1`, kind); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return output, nil
}

func (r *Place) ImagesByPlaceID(id string) ([]entity.UUID, error) {
	var output []entity.UUID

	if err := r.db.Select(&output, `SELECT image_id FROM place_image WHERE place_id = $1`, id); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return output, nil
}
