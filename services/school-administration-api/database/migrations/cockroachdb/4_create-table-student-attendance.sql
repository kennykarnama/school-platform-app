-- +goose Up
CREATE TABLE IF NOT EXISTS student_attendance (
    id UUID primary key NOT null default gen_random_uuid(),
    student_class_id UUID NOT NULL,
    attend BOOL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS student_attendance;
