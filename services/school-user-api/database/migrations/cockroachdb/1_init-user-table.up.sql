-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() primary key ,
    alternative_id STRING NOT NULL UNIQUE,
    password STRING NOT NULL,
    name STRING NOT NULL ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS users;
