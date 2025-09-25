package tour

import (
	"arch/internal/domain/entity"
	"encoding/json"
	"errors"
	"time"
)

type planEnvelopeDays struct {
	DateTour *entity.DateTour `json:"date_tour"`
	Days     []entity.DayPlan `json:"placesInfo"`
}

func normalizeTourOutput(dbFrom, dbTo time.Time, env *planEnvelopeDays, raw []byte) (entity.TourOutput, error) {
	// Базовые даты — из БД
	out := entity.TourOutput{
		DateTour: entity.DateTour{
			DateFrom: dbFrom.Format("2006-01-02"),
			DateTo:   dbTo.Format("2006-01-02"),
		},
		Days: nil,
	}

	if env != nil && env.DateTour != nil {
		if env.DateTour.DateFrom != "" {
			out.DateFrom = env.DateTour.DateFrom
		}
		if env.DateTour.DateTo != "" {
			out.DateTo = env.DateTour.DateTo
		}
	}

	if env != nil && len(env.Days) > 0 {
		out.Days = env.Days
		return out, nil
	}

	var asDays []entity.DayPlan
	if err := json.Unmarshal(raw, &asDays); err == nil && len(asDays) > 0 && asDays[0].Day != "" {
		out.Days = asDays
		return out, nil
	}

	var flat []entity.PlaceInfo
	if err := json.Unmarshal(raw, &flat); err == nil && len(flat) > 0 {
		out.Days = []entity.DayPlan{
			{Day: out.DateFrom, Places: flat},
		}
		return out, nil
	}

	return out, errors.New("unsupported plan format")
}

func (r *Tour) Reader(sessionID string) ([]entity.TourOutput, error) {
	type row struct {
		DateFrom time.Time `db:"date_from"`
		DateTo   time.Time `db:"date_to"`
		Plan     []byte    `db:"plan"`
	}

	var rows []row
	if err := r.db.Select(&rows, `
SELECT date_from, date_to, plan
FROM tour
WHERE session = $1
ORDER BY date_from DESC
`, sessionID); err != nil {
		return nil, err
	}

	out := make([]entity.TourOutput, 0, len(rows))
	for _, rr := range rows {
		var env planEnvelopeDays
		_ = json.Unmarshal(rr.Plan, &env)
		tout, _ := normalizeTourOutput(rr.DateFrom, rr.DateTo, &env, rr.Plan)
		out = append(out, tout)
	}
	return out, nil
}

func (r *Tour) ReaderByID(id entity.UUID) (*entity.TourOutput, error) {
	type row struct {
		DateFrom time.Time `db:"date_from"`
		DateTo   time.Time `db:"date_to"`
		Plan     []byte    `db:"plan"`
	}

	var rr row
	if err := r.db.Get(&rr, `
SELECT date_from, date_to, plan
FROM tour
WHERE id = $1::uuid
LIMIT 1
`, id); err != nil {
		return nil, err
	}

	var env planEnvelopeDays
	_ = json.Unmarshal(rr.Plan, &env)
	tout, _ := normalizeTourOutput(rr.DateFrom, rr.DateTo, &env, rr.Plan)
	return &tout, nil
}
