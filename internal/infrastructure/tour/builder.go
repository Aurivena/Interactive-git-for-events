package tour

import "strings"

func buildFullSQL() string {
	parts := []string{
		"WITH RECURSIVE",
		buildParamsCTE() + ",",
		buildDaysCTE() + ",",
		buildWeekdayMapCTE() + ",",
		buildCandidatesCTE() + ",",
		buildScheduleFilterCTE() + ",",
		buildRankedCTE() + ",",
		buildAssignedCTE() + ",",
		buildDayKindStatsCTE() + ",", // <— НОВОЕ
		buildPerDayCTE() + ",",
		buildStartNodesCTE() + ",",
		buildRouteCTE(),
		buildFinalSelect(),
	}
	return strings.Join(parts, "\n")
}

func buildParamsCTE() string {
	return `
cte_params AS (
  SELECT
    $1::date        AS date_from,
    $2::date        AS date_to,
    $3::float8      AS start_lon,
    $4::float8      AS start_lat,
    $5::int         AS per_day_limit,
    $6::tier_enum   AS max_tier,
    $7::kind_enum[] AS kind_priority,
    $8::time        AS day_start,
    $9::time        AS day_end
)`
}

func buildDaysCTE() string {
	return `
cte_days AS (
  SELECT generate_series((SELECT date_from FROM cte_params),
                         (SELECT date_to   FROM cte_params),
                         interval '1 day')::date AS day
)`
}

func buildWeekdayMapCTE() string {
	return `
cte_weekday_map AS (
  SELECT
    d.day,
    CASE extract(isodow from d.day)::int
      WHEN 1 THEN 'monday'
      WHEN 2 THEN 'tuesday'
      WHEN 3 THEN 'wednesday'
      WHEN 4 THEN 'thursday'
      WHEN 5 THEN 'friday'
      WHEN 6 THEN 'saturday'
      WHEN 7 THEN 'sunday'
    END AS week_txt
  FROM cte_days d
)`
}

func buildCandidatesCTE() string {
	return `
cte_candidates AS (
  SELECT
    d.day,
    pl.id, pl.title, pl.address, pl.lon, pl.lat, pl.kind, pl.tier, pl.tags,
    2 * 6371 * asin(
      sqrt(
        sin(radians((pl.lat - (SELECT start_lat FROM cte_params))/2))^2 +
        cos(radians((SELECT start_lat FROM cte_params))) * cos(radians(pl.lat)) *
        sin(radians((pl.lon - (SELECT start_lon FROM cte_params))/2))^2
      )
    ) AS dist_from_start,
    COALESCE(array_position((SELECT kind_priority FROM cte_params), pl.kind), 999) AS kind_rank
  FROM cte_days d
  JOIN place pl
    ON pl.tier <= (SELECT max_tier FROM cte_params)
   AND (SELECT kind_priority FROM cte_params) @> ARRAY[pl.kind]::kind_enum[]
)`
}

func buildScheduleFilterCTE() string {
	return `
cte_open_candidates AS (
  SELECT c.*
  FROM cte_candidates c
  JOIN cte_weekday_map wm ON wm.day = c.day
  WHERE
    c.tags->'schedule' IS NULL
    OR EXISTS (
      SELECT 1
      FROM jsonb_array_elements(c.tags->'schedule') s(obj)
      WHERE
        lower(trim((obj->>'week'))) = wm.week_txt
        AND (
          (
            COALESCE((obj->>'spans_midnight')::boolean, false) = false
            AND GREATEST((SELECT day_start FROM cte_params),(obj->>'start')::time)
                < LEAST((SELECT day_end FROM cte_params),(obj->>'end')::time)
          )
          OR
          (
            COALESCE((obj->>'spans_midnight')::boolean, false) = true
            AND ( (SELECT day_end FROM cte_params) > (obj->>'start')::time
               OR (SELECT day_start FROM cte_params) < (obj->>'end')::time )
          )
        )
    )
)`
}

func buildRankedCTE() string {
	return `
cte_unique_open AS (
  SELECT *
  FROM (
    SELECT
      c.*,
      row_number() OVER (
        PARTITION BY c.id, c.day
        ORDER BY c.kind_rank NULLS LAST, c.dist_from_start
      ) AS day_rn
    FROM cte_open_candidates c
  ) t
  WHERE day_rn = 1
),
cte_ranked AS (
  SELECT *,
         row_number() OVER (ORDER BY kind_rank NULLS LAST, dist_from_start, id) AS global_rank
  FROM cte_unique_open
  LIMIT 300
)`
}

func buildAssignedCTE() string {
	return `
cte_assigned AS (
  SELECT r.*
  FROM cte_ranked r
)`
}

func buildDayKindStatsCTE() string {
	return `
cte_day_kind_stats AS (
  SELECT day, count(DISTINCT kind) AS kinds_cnt
  FROM cte_ranked
  GROUP BY day
)`
}

func buildPerDayCTE() string {
	return `
cte_per_day AS (
  SELECT
    a.day,
    a.id, a.title, a.address, a.lon, a.lat, a.kind, a.tier, a.tags,
    a.dist_from_start, a.kind_rank,
    row_number() OVER (
      PARTITION BY a.day
      ORDER BY a.kind_rank NULLS LAST, a.dist_from_start, a.id
    ) AS slot_rank,
    row_number() OVER (
      PARTITION BY a.day, a.kind
      ORDER BY a.kind_rank NULLS LAST, a.dist_from_start, a.id
    ) AS per_kind_rn,
    GREATEST(
      1,
      CEIL(
        (SELECT per_day_limit FROM cte_params)::numeric
        / NULLIF((SELECT kinds_cnt FROM cte_day_kind_stats dks WHERE dks.day = a.day), 0)
      )::numeric
    )::int AS per_kind_cap
  FROM cte_ranked a
)`
}

func buildStartNodesCTE() string {
	return `
cte_start_nodes AS (
  SELECT
    pd.day,
    1 AS step,
    NULL::uuid AS prev_id,
    (SELECT start_lon FROM cte_params) AS cur_lon,
    (SELECT start_lat FROM cte_params) AS cur_lat,
    ARRAY[]::uuid[] AS visited,
    ARRAY[]::kind_enum[] AS visited_kinds
  FROM (SELECT DISTINCT day FROM cte_per_day) pd
)`
}
func buildRouteCTE() string {
	return `
cte_route AS (
  -- старт на каждый день: берём ближайшую точку, но с приоритетом "ещё не было такого вида"
  SELECT
    s.day,
    1 AS step,
    p.id, p.title, p.address, p.lon, p.lat, p.kind, p.tier, p.tags,
    2 * 6371 * asin(
      sqrt(
        sin(radians((p.lat - s.cur_lat)/2))^2 +
        cos(radians(s.cur_lat)) * cos(radians(p.lat)) *
        sin(radians((p.lon - s.cur_lon)/2))^2
      )
    ) AS leg_km,
    ARRAY[p.id] AS visited,
    ARRAY[p.kind]::kind_enum[] AS visited_kinds
  FROM cte_start_nodes s
  JOIN LATERAL (
    SELECT * FROM cte_per_day pd
    WHERE
      pd.day = s.day
      AND pd.slot_rank   <= GREATEST(1, (SELECT per_day_limit FROM cte_params))
      AND pd.per_kind_rn <= pd.per_kind_cap
ORDER BY
  pd.per_kind_rn,
  2 * 6371 * asin(
    sqrt(
      sin(radians((pd.lat - s.cur_lat)/2))^2 +
      cos(radians(s.cur_lat)) * cos(radians(pd.lat)) *
      sin(radians((pd.lon - s.cur_lon)/2))^2
    )
  )
    LIMIT 1
  ) p ON TRUE

  UNION ALL

  -- продолжение: сначала вид, которого ещё не было в этом дне, и не превышаем cap
  SELECT
    r.day,
    r.step + 1,
    p.id, p.title, p.address, p.lon, p.lat, p.kind, p.tier, p.tags,
    2 * 6371 * asin(
      sqrt(
        sin(radians((p.lat - r.lat)/2))^2 +
        cos(radians(r.lat)) * cos(radians(p.lat)) *
        sin(radians((p.lon - r.lon)/2))^2
      )
    ) AS leg_km,
    r.visited || p.id,
    r.visited_kinds || p.kind
  FROM cte_route r
  JOIN LATERAL (
    SELECT * FROM cte_per_day pd
    WHERE
      pd.day = r.day
      AND NOT (pd.id = ANY (r.visited))
      AND pd.slot_rank   <= GREATEST(1, (SELECT per_day_limit FROM cte_params))
      AND pd.per_kind_rn <= pd.per_kind_cap
ORDER BY
  CASE WHEN NOT (pd.kind = ANY (r.visited_kinds)) THEN 0 ELSE 1 END,
  pd.per_kind_rn,
  2 * 6371 * asin(
    sqrt(
      sin(radians((pd.lat - r.lat)/2))^2 +
      cos(radians(r.lat)) * cos(radians(pd.lat)) *
      sin(radians((pd.lon - r.lon)/2))^2
    )
  )
    LIMIT 1
  ) p ON TRUE
  WHERE r.step < GREATEST(1, (SELECT per_day_limit FROM cte_params))
)`
}

func buildFinalSelect() string {
	return `
SELECT jsonb_build_object(
  'date_tour', jsonb_build_object(
      'date_from', (SELECT date_from FROM cte_params),
      'date_to',   (SELECT date_to   FROM cte_params)
  ),
  'placesInfo', COALESCE(jsonb_agg(day_block ORDER BY day), '[]'::jsonb)
)
FROM (
  SELECT
    day,
    jsonb_build_object(
      'day', day,
      'places', jsonb_agg(
        jsonb_build_object(
          'step', step,
          'id', id,
          'title', title,
          'address', address,
          'lon', lon,
          'lat', lat,
          'kind', kind::text,
          'tier', tier::text,
          'leg_km', round(leg_km::numeric, 2)
        ) ORDER BY step
      )
    ) AS day_block
  FROM cte_route
  GROUP BY day
) t`
}
