package tour

import (
	"arch/internal/domain/entity"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func (r *Tour) GenerateTour(input entity.RouteParams, lon, lat float64) (json.RawMessage, error) {

	var route json.RawMessage
	if err := r.db.Get(&route, buildFullSQL(), input.DateFrom, input.DateTo, lon, lat, input.PerDayLimit, input.Tier, pq.Array(input.KindPriority), input.DayStart, input.DayEnd); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return route, nil
}
