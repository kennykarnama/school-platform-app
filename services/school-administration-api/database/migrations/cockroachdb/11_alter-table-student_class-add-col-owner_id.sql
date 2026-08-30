-- +goose Up
ALTER TABLE student_class ADD COLUMN IF NOT EXISTS owner_id STRING;

-- +goose Down
ALTER TABLE student_class DROP COLUMN IF EXISTS owner_id;
