-- +goose Up
CREATE TABLE IF NOT EXISTS teacher (
    id UUID primary key NOT null default gen_random_uuid(),
    alternative_id uuid NOT NULL,
    name STRING NOT NULL,
    password STRING NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS teacher;
