-- +goose Up
ALTER TABLE attendance_type ADD COLUMN IF NOT EXISTS color STRING;

-- +goose Down
ALTER TABLE attendance_type DROP COLUMN IF EXISTS color;
