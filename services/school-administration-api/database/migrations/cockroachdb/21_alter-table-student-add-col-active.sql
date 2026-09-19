-- +goose Up
ALTER TABLE student ADD COLUMN IF NOT EXISTS active BOOL NOT NULL DEFAULT TRUE;
ALTER TABLE student ADD COLUMN IF NOT EXISTS deactivate_reason STRING;
CREATE INDEX IF NOT EXISTS student_school_active_idx ON student (school_id, active);

-- +goose Down
DROP INDEX IF EXISTS student_school_active_idx;
ALTER TABLE student DROP COLUMN IF EXISTS deactivate_reason;
ALTER TABLE student DROP COLUMN IF EXISTS active;
