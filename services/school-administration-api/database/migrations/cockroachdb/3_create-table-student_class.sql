-- +goose Up
CREATE TABLE IF NOT EXISTS student_class (
    id UUID primary key NOT null default gen_random_uuid(),
    student_id UUID NOT NULL,
    class_label STRING NOT NULL,
    academic_year_id UUID NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS student_class;
