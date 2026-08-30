-- +goose Up
ALTER TABLE student_class ADD COLUMN IF NOT EXISTS deactivate_reason STRING;

-- +goose Down
ALTER TABLE student_class DROP COLUMN IF EXISTS deactivate_reason;
