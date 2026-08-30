-- +goose Up
CREATE TABLE IF NOT EXISTS student (
    id UUID primary key NOT null default gen_random_uuid(),
    name STRING NOT NULL,
    alternative_id STRING,
    graduated BOOL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS student;
