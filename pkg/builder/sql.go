package builder

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

func BuildSql(params *entity.RequestPayload, base string, centerLon, centerLat *float64) (string, []any) {
	clauses := make([]string, 0, 8)
	args := make([]any, 0, 8)

	addArg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	if params.Kind != "" {
		clauses = append(clauses, fmt.Sprintf("p.kind = %s::kind_enum", addArg(params.Kind)))
	}

	if params.Tier != "" {
		clauses = append(clauses, fmt.Sprintf("p.tier = %s::tier_enum", addArg(params.Tier)))
	}

	needSchedule := false
	schedConds := make([]string, 0, 3)

	if len(params.DayOfTheWeek) > 0 {
		weeks := make([]string, 0, len(params.DayOfTheWeek))
		for _, w := range params.DayOfTheWeek {
			if w.Valid() {
				weeks = append(weeks, w.Convert()) // "monday" и т.п.
			}
		}
		if len(weeks) > 0 {
			needSchedule = true
			schedConds = append(schedConds, fmt.Sprintf("(s.week = ANY(%s::text[]))", addArg(pq.Array(weeks))))
		}
	}

	// "открыто в момент t"
	if params.Time != nil && params.Time.Valid() {
		needSchedule = true
		t := fmt.Sprintf("%02d:%02d", params.Time.Hour, params.Time.Minute)
		ti := addArg(t)

		// spans_midnight как boolean (COALESCE на случай отсутствия поля)
		schedConds = append(schedConds, fmt.Sprintf(`
(
  (COALESCE(s.spans_midnight, false) = false AND %s::time BETWEEN (s.start)::time AND (s."end")::time)
  OR
  (COALESCE(s.spans_midnight, false) = true  AND (%s::time >= (s.start)::time OR %s::time < (s."end")::time))
)`, ti, ti, ti))
	}

	if needSchedule {
		// Берём массив schedule из p.tags и разворачиваем в строки s(...)
		clauses = append(clauses, fmt.Sprintf(`
EXISTS (
  SELECT 1
  FROM jsonb_to_recordset(p.tags->'schedule')
       AS s(week text, start text, "end" text, spans_midnight boolean)
  WHERE %s
)`, strings.Join(schedConds, " AND ")))
	}
	// --- КОНЕЦ ЗАМЕНЫ ---

	// Радиус — без изменений
	if params.Radius > 0 && centerLon != nil && centerLat != nil {
		clauses = append(clauses,
			fmt.Sprintf(`ST_DWithin(p.geom, ST_MakePoint(%s,%s)::geography, %s)`,
				addArg(*centerLon), addArg(*centerLat), addArg(params.Radius)))
	}

	// Склейка
	query := base
	if len(clauses) > 0 {
		if !strings.Contains(strings.ToUpper(base), "WHERE") {
			query += " WHERE 1=1"
		}
		query += " AND " + strings.Join(clauses, " AND ")
	}

	return query, args
}
