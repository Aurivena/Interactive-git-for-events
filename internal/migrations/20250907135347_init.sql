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

CREATE TABLE IF NOT EXISTS place_image(
    place_id uuid,
    image_id uuid,
    PRIMARY KEY (place_id,image_id)
);

CREATE TABLE IF NOT EXISTS history (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session varchar(255) NOT NULL ,
    message text NOT NULL ,
    ai_message jsonb NOT NULL,
    created_at timestamp DEFAULT NOW() NOT NULL
);

ALTER TABLE place_image
    ADD CONSTRAINT fk_place_image_place
        FOREIGN KEY (place_id) REFERENCES place(id) ON DELETE CASCADE;


CREATE INDEX IF NOT EXISTS idx_tier_kind ON place(tier,kind);
CREATE INDEX IF NOT EXISTS idx_place_trgm ON place USING gin((title || ' ' || address) gin_trgm_ops);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_place_trgm;
DROP INDEX IF EXISTS idx_tier_kind;
DROP INDEX IF EXISTS idx_place_image_place;

DROP TABLE IF EXISTS history;
DROP TABLE IF EXISTS place_image;
DROP TABLE IF EXISTS place;

DROP TYPE IF EXISTS kind_enum;
DROP TYPE IF EXISTS tier_enum;
DROP TYPE IF EXISTS week_enum;
-- +goose StatementEnd
