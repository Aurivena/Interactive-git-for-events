-- +goose Up
-- +goose StatementBegin
-- schema
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE  week_enum as enum('monday','tuesday','wednesday','thursday','friday','saturday','sunday');
CREATE TYPE tier_enum as enum('economy','value','standard','premium','upscale');
CREATE TYPE kind_enum  AS ENUM (
    'cinema','theatre','concert_hall','stadium','sport','museum','art_gallery',
    'historic','memorial','park','zoo','aquapark','attraction',
    'church','monastery','mosque','synagogue','mall','market','monument','restaurant'
    );

CREATE TABLE IF NOT EXISTS place (
    id       uuid  PRIMARY KEY,
    title    varchar(300) NOT NULL,
    address  varchar(300) not null ,
    description text,
    lon      double precision NOT NULL ,
    lat      double precision NOT NULL ,
    tier tier_enum NOT NULL ,
    kind     kind_enum NOT NULL,
    tags jsonb
);

CREATE TABLE IF NOT EXISTS place_schedule (
    place_id uuid NOT NULL,
    start_work time,
    end_work time,
    week week_enum NOT NULL,
    spans_midnight boolean NOT NULL DEFAULT false

        CONSTRAINT chk_times_null_pair CHECK (
            (start_work IS NULL AND end_work IS NULL)
                OR (start_work IS NOT NULL AND end_work IS NOT NULL)
            ),
    CONSTRAINT chk_times_order CHECK (
        (start_work IS NULL AND end_work IS NULL AND spans_midnight = false)
            OR (start_work IS NOT NULL AND end_work IS NOT NULL AND (
            (spans_midnight = false AND start_work < end_work)
                OR (spans_midnight = true)
            ))
        )
);

CREATE TABLE IF NOT EXISTS history (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session varchar(255) NOT NULL ,
    message text NOT NULL ,
    ai_message jsonb NOT NULL,
    created_at timestamp DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tier_kind ON place(tier,kind);
CREATE INDEX IF NOT EXISTS idx_place_trgm ON place USING gin((title || ' ' || address) gin_trgm_ops);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS place;
DROP TABLE IF EXISTS place_schedule;
-- +goose StatementEnd
