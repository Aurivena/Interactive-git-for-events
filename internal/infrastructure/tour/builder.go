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
		buildDayKindStatsCTE() + ",",
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
    $1::date                  AS date_from,
    $2::date                  AS date_to,
    $3::float8                AS start_lon,
    $4::float8                AS start_lat,
    $5::int                   AS per_day_limit,
    ($6::text[])::tier_enum[] AS allowed_tiers,
    ($7::text[])::kind_enum[] AS kind_priority,
    $8::time                  AS day_start,
    $9::time                  AS day_end
)`
}
func buildDaysCTE() string {
	return `
cte_days AS (
  SELECT generate_series(
           (SELECT date_from FROM cte_params),
           (SELECT date_to   FROM cte_params),
           interval '1 day'
         )::date AS day
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
    pl.id, pl.title, pl.address, pl.description, pl.lon, pl.lat, pl.kind, pl.tier, pl.tags,
    img.images,
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
    ON (
         (SELECT COALESCE(cardinality(allowed_tiers),0) FROM cte_params) = 0
         OR EXISTS (
              SELECT 1
              FROM unnest((SELECT allowed_tiers FROM cte_params)) AS t(v)
              WHERE t.v = pl.tier
            )
       )
  LEFT JOIN LATERAL (
    SELECT COALESCE(array_agg(pi.image_id ORDER BY pi.image_id), ARRAY[]::uuid[]) AS images
    FROM place_image pi
    WHERE pi.place_id = pl.id
  ) img ON TRUE
)`
}

func buildScheduleFilterCTE() string {
	return `
cte_open_candidates AS (
  WITH norm AS (
    SELECT c.*, wm.week_txt AS needed_day
    FROM cte_candidates c
    JOIN cte_weekday_map wm ON wm.day = c.day
  )
  SELECT n.*
  FROM norm n
  WHERE
      n.tags->'schedule' IS NULL
      OR jsonb_typeof(n.tags->'schedule') <> 'array'
      OR jsonb_array_length(n.tags->'schedule') = 0
      OR EXISTS (
        SELECT 1
        FROM jsonb_array_elements(n.tags->'schedule') s(obj)
        WHERE
          CASE lower(trim(obj->>'week'))
            WHEN 'monday'      THEN 'monday'
            WHEN 'tuesday'     THEN 'tuesday'
            WHEN 'wednesday'   THEN 'wednesday'
            WHEN 'thursday'    THEN 'thursday'
            WHEN 'friday'      THEN 'friday'
            WHEN 'saturday'    THEN 'saturday'
            WHEN 'sunday'      THEN 'sunday'
            WHEN 'понедельник' THEN 'monday'
            WHEN 'вторник'     THEN 'tuesday'
            WHEN 'среда'       THEN 'wednesday'
            WHEN 'четверг'     THEN 'thursday'
            WHEN 'пятница'     THEN 'friday'
            WHEN 'суббота'     THEN 'saturday'
            WHEN 'воскресенье' THEN 'sunday'
            ELSE NULL
          END = n.needed_day
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
-- 1) Дедуп внутри дня (одно место один раз в конкретный день)
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

-- 2) Индексация дней периода
cte_days_idx AS (
  SELECT
    d.day,
    row_number() OVER (ORDER BY d.day) AS day_idx,
    count(*)    OVER ()                AS days_total
  FROM cte_days d
),

-- 3) Стабильный порядок мест (один раз на id)
cte_place_order AS (
  SELECT
    uo.id,
    min(uo.kind_rank)       AS min_kind_rank,
    min(uo.dist_from_start) AS min_dist
  FROM cte_unique_open uo
  GROUP BY uo.id
),

-- 4) Раскладываем все (id, day) по сетке
cte_spread AS (
  SELECT
    uo.*,
    di.day_idx,
    di.days_total,
    row_number() OVER (PARTITION BY uo.id ORDER BY di.day_idx) AS day_rank_for_place,
    count(*)    OVER (PARTITION BY uo.id)                       AS place_days_total,
    row_number() OVER (
      ORDER BY po.min_kind_rank NULLS LAST, po.min_dist, uo.id
    )                                                           AS place_idx
  FROM cte_unique_open uo
  JOIN cte_days_idx   di ON di.day = uo.day
  JOIN cte_place_order po ON po.id  = uo.id
),

-- 5) Выбираем РОВНО один день на каждое place.id (равномерно)
cte_picked AS (
  SELECT *
  FROM cte_spread
  WHERE day_rank_for_place = ((place_idx - 1) % place_days_total) + 1
),

-- 6) На всякий пожарный убираем любые остаточные дубли по id
cte_picked_distinct AS (
  SELECT DISTINCT ON (id)
         *
  FROM cte_picked
  ORDER BY id, day_idx, kind_rank NULLS LAST, dist_from_start
),

-- 7) Финальный пул
cte_ranked AS (
  SELECT *,
         row_number() OVER (ORDER BY kind_rank NULLS LAST, dist_from_start, id) AS global_rank
  FROM cte_picked_distinct
  LIMIT 10000
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
    a.id, a.title, a.address, a.description, a.lon, a.lat, a.kind, a.tier, a.tags, a.images,
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
      LEAST(
        (SELECT per_day_limit FROM cte_params),
        CEIL(
          (SELECT per_day_limit FROM cte_params)::numeric
          / NULLIF((SELECT kinds_cnt FROM cte_day_kind_stats dks WHERE dks.day = a.day), 0)
        )::numeric
      )
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
cte_route (
  day, step,
  id, title, address, description, lon, lat, kind, tier, tags, images,
  leg_km, visited, visited_kinds
) AS (
  SELECT
    s.day,
    1 AS step,
    c.id, c.title, c.address, c.description, c.lon, c.lat, c.kind, c.tier, c.tags, c.images,
    2 * 6371 * asin(
      sqrt(
        sin(radians((c.lat - s.cur_lat)/2))^2 +
        cos(radians(s.cur_lat)) * cos(radians(c.lat)) *
        sin(radians((c.lon - s.cur_lon)/2))^2
      )
    ) AS leg_km,
    ARRAY[c.id]                AS visited,
    ARRAY[c.kind]::kind_enum[] AS visited_kinds
  FROM cte_start_nodes s
  CROSS JOIN LATERAL (
    SELECT pd.*
    FROM cte_per_day pd
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
  ) AS c

  UNION ALL

  SELECT
    r.day,
    r.step + 1,
    c.id, c.title, c.address, c.description, c.lon, c.lat, c.kind, c.tier, c.tags, c.images,
    2 * 6371 * asin(
      sqrt(
        sin(radians((c.lat - r.lat)/2))^2 +
        cos(radians(r.lat)) * cos(radians(c.lat)) *
        sin(radians((c.lon - r.lon)/2))^2
      )
    ) AS leg_km,
    r.visited || c.id,
    r.visited_kinds || c.kind
  FROM cte_route r
  CROSS JOIN LATERAL (
    SELECT pd.*
    FROM cte_per_day pd
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
  ) AS c
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
    d.day,
    jsonb_build_object(
      'day', d.day,
      'places',
        COALESCE(
          (
            SELECT jsonb_agg(
                     jsonb_build_object(
                       'step', r.step,
                       'id', r.id,
                       'title', r.title,
                       'description', r.description,
                       'address', r.address,
                       'lon', r.lon,
                       'lat', r.lat,
                       'kind', r.kind::text,
                       'tier', r.tier::text,
                       'tags', COALESCE(r.tags, '{}'::jsonb),
                       'images', to_jsonb(COALESCE(r.images, ARRAY[]::uuid[])),
                       'leg_km', round(r.leg_km::numeric, 2)
                     )
                     ORDER BY r.step
                   )
            FROM cte_route r
            WHERE r.day = d.day
          ),
          '[]'::jsonb
        )
    ) AS day_block
  FROM cte_days d
) t`
}
