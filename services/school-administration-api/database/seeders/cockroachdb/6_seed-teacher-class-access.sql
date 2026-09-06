INSERT INTO teacher_class_access (
    id,
    school_id,
    teacher_id,
    academic_year_id,
    class_label,
    created_at
) VALUES (
    '70000000-0000-4000-8000-000000000001',
    '00000000-0000-4000-8000-000000000001',
    '30000000-0000-4000-8000-000000000001',
    '10000000-0000-4000-8000-000000000002',
    'KELAS I A',
    '2026-07-13T00:00:00Z'
)
ON CONFLICT (school_id, teacher_id, academic_year_id, class_label) DO NOTHING;
