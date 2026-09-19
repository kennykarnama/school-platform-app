-- +goose Up
CREATE TABLE IF NOT EXISTS school (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    code STRING NOT NULL,
    active BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT school_code_unique UNIQUE (code)
);

INSERT INTO school (id, name, code)
VALUES ('00000000-0000-4000-8000-000000000001', 'Existing School', 'legacy')
ON CONFLICT (code) DO NOTHING;

ALTER TABLE teacher ADD COLUMN IF NOT EXISTS school_id UUID;
ALTER TABLE teacher ADD COLUMN IF NOT EXISTS role STRING NOT NULL DEFAULT 'teacher';
ALTER TABLE teacher ADD COLUMN IF NOT EXISTS active BOOL NOT NULL DEFAULT TRUE;
ALTER TABLE teacher ADD COLUMN IF NOT EXISTS must_change_password BOOL NOT NULL DEFAULT FALSE;
UPDATE teacher SET school_id = '00000000-0000-4000-8000-000000000001' WHERE school_id IS NULL;
UPDATE teacher SET role = 'school_admin'
WHERE id = (SELECT id FROM teacher WHERE deleted_at IS NULL ORDER BY created_at ASC, id ASC LIMIT 1);
ALTER TABLE teacher ADD CONSTRAINT teacher_role_valid CHECK (role IN ('platform_admin', 'school_admin', 'teacher'));
ALTER TABLE teacher ADD CONSTRAINT teacher_school_role_valid CHECK (
    (role = 'platform_admin' AND school_id IS NULL) OR
    (role IN ('school_admin', 'teacher') AND school_id IS NOT NULL)
);
CREATE UNIQUE INDEX IF NOT EXISTS teacher_alternative_id_unique ON teacher (alternative_id);

ALTER TABLE academic_year ADD COLUMN IF NOT EXISTS school_id UUID;
UPDATE academic_year SET school_id = '00000000-0000-4000-8000-000000000001' WHERE school_id IS NULL;
ALTER TABLE academic_year ALTER COLUMN school_id SET NOT NULL;

ALTER TABLE attendance_type ADD COLUMN IF NOT EXISTS school_id UUID;
UPDATE attendance_type SET school_id = '00000000-0000-4000-8000-000000000001' WHERE school_id IS NULL;
ALTER TABLE attendance_type ALTER COLUMN school_id SET NOT NULL;

ALTER TABLE student ADD COLUMN IF NOT EXISTS school_id UUID;
UPDATE student SET school_id = '00000000-0000-4000-8000-000000000001' WHERE school_id IS NULL;
ALTER TABLE student ALTER COLUMN school_id SET NOT NULL;
DROP INDEX IF EXISTS student_alternative_id_unique;
CREATE UNIQUE INDEX IF NOT EXISTS student_school_alternative_id_unique ON student (school_id, alternative_id);

ALTER TABLE student_class ADD COLUMN IF NOT EXISTS school_id UUID;
UPDATE student_class SET school_id = '00000000-0000-4000-8000-000000000001' WHERE school_id IS NULL;
ALTER TABLE student_class ALTER COLUMN school_id SET NOT NULL;

ALTER TABLE student_attendance ADD COLUMN IF NOT EXISTS school_id UUID;
UPDATE student_attendance AS attendance
SET school_id = student_class.school_id
FROM student_class
WHERE attendance.student_class_id = student_class.id AND attendance.school_id IS NULL;
UPDATE student_attendance SET school_id = '00000000-0000-4000-8000-000000000001' WHERE school_id IS NULL;
ALTER TABLE student_attendance ALTER COLUMN school_id SET NOT NULL;

CREATE TABLE IF NOT EXISTS school_class (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL,
    label STRING NOT NULL,
    active BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT school_class_unique UNIQUE (school_id, label)
);

INSERT INTO school_class (school_id, label)
SELECT '00000000-0000-4000-8000-000000000001', label
FROM (VALUES
    ('KELAS I A'), ('KELAS I B'), ('KELAS I C'),
    ('KELAS II A'), ('KELAS II B'), ('KELAS II C'),
    ('KELAS III A'), ('KELAS III B'), ('KELAS III C'),
    ('KELAS IV A'), ('KELAS IV B'), ('KELAS IV C'),
    ('KELAS V A'), ('KELAS V B'), ('KELAS V C'),
    ('KELAS VI A'), ('KELAS VI B'), ('KELAS VI C')
) AS defaults(label)
ON CONFLICT (school_id, label) DO NOTHING;

INSERT INTO school_class (school_id, label)
SELECT DISTINCT school_id, class_label FROM student_class WHERE trim(class_label) <> ''
ON CONFLICT (school_id, label) DO NOTHING;

CREATE TABLE IF NOT EXISTS teacher_class_access (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL,
    teacher_id UUID NOT NULL,
    academic_year_id UUID NOT NULL,
    class_label STRING NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT teacher_class_access_unique UNIQUE (school_id, teacher_id, academic_year_id, class_label)
);

INSERT INTO teacher_class_access (school_id, teacher_id, academic_year_id, class_label)
SELECT DISTINCT student_class.school_id, teacher.id, student_class.academic_year_id, student_class.class_label
FROM student_class
JOIN teacher ON teacher.id::STRING = student_class.user_id
WHERE student_class.user_id IS NOT NULL AND trim(student_class.user_id) <> ''
ON CONFLICT (school_id, teacher_id, academic_year_id, class_label) DO NOTHING;

CREATE INDEX IF NOT EXISTS academic_year_school_idx ON academic_year (school_id);
CREATE INDEX IF NOT EXISTS attendance_type_school_idx ON attendance_type (school_id);
CREATE INDEX IF NOT EXISTS student_class_school_idx ON student_class (school_id, academic_year_id, class_label);
CREATE INDEX IF NOT EXISTS student_attendance_school_idx ON student_attendance (school_id);
CREATE INDEX IF NOT EXISTS teacher_school_idx ON teacher (school_id, role, active);
CREATE INDEX IF NOT EXISTS teacher_class_access_lookup_idx ON teacher_class_access (teacher_id, academic_year_id, class_label);

-- +goose Down
DROP INDEX IF EXISTS teacher_class_access_lookup_idx;
DROP INDEX IF EXISTS teacher_school_idx;
DROP INDEX IF EXISTS student_attendance_school_idx;
DROP INDEX IF EXISTS student_class_school_idx;
DROP INDEX IF EXISTS attendance_type_school_idx;
DROP INDEX IF EXISTS academic_year_school_idx;
DROP TABLE IF EXISTS teacher_class_access;
DROP TABLE IF EXISTS school_class;
DROP INDEX IF EXISTS student_school_alternative_id_unique;
CREATE UNIQUE INDEX IF NOT EXISTS student_alternative_id_unique ON student (alternative_id);
ALTER TABLE student_attendance DROP COLUMN IF EXISTS school_id;
ALTER TABLE student_class DROP COLUMN IF EXISTS school_id;
ALTER TABLE student DROP COLUMN IF EXISTS school_id;
ALTER TABLE attendance_type DROP COLUMN IF EXISTS school_id;
ALTER TABLE academic_year DROP COLUMN IF EXISTS school_id;
DROP INDEX IF EXISTS teacher_alternative_id_unique;
ALTER TABLE teacher DROP CONSTRAINT IF EXISTS teacher_school_role_valid;
ALTER TABLE teacher DROP CONSTRAINT IF EXISTS teacher_role_valid;
ALTER TABLE teacher DROP COLUMN IF EXISTS must_change_password;
ALTER TABLE teacher DROP COLUMN IF EXISTS active;
ALTER TABLE teacher DROP COLUMN IF EXISTS role;
ALTER TABLE teacher DROP COLUMN IF EXISTS school_id;
DROP TABLE IF EXISTS school;
