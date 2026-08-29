CREATE TABLE attendance_type (
    id UUID primary key NOT null default gen_random_uuid(),
    label STRING NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);