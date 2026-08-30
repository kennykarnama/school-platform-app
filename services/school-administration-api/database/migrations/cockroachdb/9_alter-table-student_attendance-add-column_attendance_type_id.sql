-- +goose Up
ALTER TABLE student_attendance ADD COLUMN IF NOT EXISTS attendance_type_id UUID;

-- +goose Down
ALTER TABLE student_attendance DROP COLUMN IF EXISTS attendance_type_id;
