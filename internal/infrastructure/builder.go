package infrastructure

import (
	"arch/internal/domain/entity"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// Пример вызова:
//   base := `SELECT p.* FROM place p
//            LEFT JOIN place_schedule ps ON ps.place_id = p.id
//            WHERE 1=1`
//   sql, args := builderSQL(&params, base, &lon, &lat)
//   // добавь GROUP BY/LIMIT по вкусу

func builderSQL(params *entity.RequestPayload, base string, centerLon, centerLat *float64) (string, []any) {
	clauses := make([]string, 0, 8)
	args := make([]any, 0, 8)

	addArg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	// kind
	if params.Kind != "" {
		clauses = append(clauses, fmt.Sprintf("p.kind = %s::kind_enum", addArg(params.Kind)))
	}

	// tier
	if params.Tier != "" {
		clauses = append(clauses, fmt.Sprintf("p.tier = %s::tier_enum", addArg(params.Tier)))
	}

	// week (массив). Явно кастуем к enum[].
	if len(params.DayOfTheWeek) > 0 {
		weeks := make([]string, 0, len(params.DayOfTheWeek))
		for _, w := range params.DayOfTheWeek {
			if w.Valid() {
				weeks = append(weeks, w.Convert())
			}
		}
		clauses = append(clauses, fmt.Sprintf("ps.week = ANY(%s::week_enum[])", addArg(pq.Array(weeks))))
	}

	// time ("открыто в этот момент")
	if params.Time != nil && params.Time.Valid() {
		t := fmt.Sprintf("%02d:%02d", params.Time.Hour, params.Time.Minute)
		ti := addArg(t) // используем один и тот же placeholder 3 раза — норм для PG

		clauses = append(clauses, fmt.Sprintf(`
(
  ps.start_work IS NOT NULL AND ps.end_work IS NOT NULL
  AND (
       (ps.spans_midnight = FALSE AND %s::time BETWEEN ps.start_work AND ps.end_work)
    OR (ps.spans_midnight = TRUE  AND (%s::time >= ps.start_work OR %s::time < ps.end_work))
  )
)`, ti, ti, ti))
	}

	// radius (PostGIS). Если нет geom — выпили или сделай bbox.
	if params.Radius > 0 && centerLon != nil && centerLat != nil {
		clauses = append(clauses,
			fmt.Sprintf(`ST_DWithin(p.geom, ST_MakePoint(%s,%s)::geography, %s)`,
				addArg(*centerLon), addArg(*centerLat), addArg(params.Radius)))
	}

	// Склейка
	query := base
	if len(clauses) > 0 {
		// base должен уже содержать WHERE 1=1 (или своё WHERE)
		if !strings.Contains(strings.ToUpper(base), "WHERE") {
			query += " WHERE 1=1"
		}
		query += " AND " + strings.Join(clauses, " AND ")
	}

	return query, args
}
