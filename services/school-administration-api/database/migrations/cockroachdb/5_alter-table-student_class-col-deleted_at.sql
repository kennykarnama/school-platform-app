-- +goose Up
ALTER TABLE student_class ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE student_class DROP COLUMN IF EXISTS deleted_at;
