-- +goose Up
ALTER TABLE student_class ADD COLUMN IF NOT EXISTS user_id STRING DEFAULT '';

-- +goose Down
ALTER TABLE student_class DROP COLUMN IF EXISTS user_id;
