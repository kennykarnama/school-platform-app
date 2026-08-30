-- +goose Up
CREATE TABLE IF NOT EXISTS academic_year (
    id UUID primary key NOT null default gen_random_uuid(),
    label STRING NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS academic_year;
