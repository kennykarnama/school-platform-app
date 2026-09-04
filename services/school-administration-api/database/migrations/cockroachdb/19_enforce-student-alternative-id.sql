-- +goose Up
UPDATE student
SET alternative_id = 'STU-' || upper(substr(replace(gen_random_uuid()::STRING, '-', ''), 1, 16))
WHERE alternative_id IS NULL OR trim(alternative_id) = '';

UPDATE student
SET alternative_id = upper(trim(alternative_id));

WITH ranked_students AS (
    SELECT
        id,
        row_number() OVER (PARTITION BY alternative_id ORDER BY created_at ASC, id ASC) AS duplicate_rank
    FROM student
)
UPDATE student
SET alternative_id = 'STU-' || upper(substr(replace(gen_random_uuid()::STRING, '-', ''), 1, 16))
FROM ranked_students
WHERE student.id = ranked_students.id
  AND ranked_students.duplicate_rank > 1;

ALTER TABLE student ALTER COLUMN alternative_id SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS student_alternative_id_unique ON student (alternative_id);

-- +goose Down
DROP INDEX IF EXISTS student_alternative_id_unique;
ALTER TABLE student ALTER COLUMN alternative_id DROP NOT NULL;
