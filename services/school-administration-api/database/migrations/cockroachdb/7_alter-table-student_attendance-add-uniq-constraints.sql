-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS class__attendance_date__unique
    ON student_attendance (student_class_id, attendance_date);

-- +goose Down
DROP INDEX IF EXISTS class__attendance_date__unique;
