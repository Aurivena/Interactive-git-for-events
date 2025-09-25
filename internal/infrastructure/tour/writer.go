package tour

import (
	"arch/internal/domain/entity"
	"encoding/json"
	"time"
)

func (r *Tour) Write(dateFrom, dateTo, sessionID string, places entity.TourOutput) (*entity.UUID, error) {

	planJSON, err := json.Marshal(places)
	if err != nil {
		return nil, err
	}

	df, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, err
	}
	dt, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return nil, err
	}

	var id entity.UUID
	if err := r.db.Get(&id, `INSERT INTO tour (session, date_from, date_to, plan)VALUES ($1, $2, $3, $4)
RETURNING id;`, sessionID, df, dt, planJSON); err != nil {
		return nil, err
	}

	return &id, nil
}
