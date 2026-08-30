-- +goose Up
ALTER TABLE academic_year ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE academic_year DROP COLUMN IF EXISTS deleted_at;
