-- Local development login: teacher.demo / school123
INSERT INTO teacher (id, school_id, alternative_id, name, password, role, active, created_at, updated_at) VALUES
    (
        '30000000-0000-4000-8000-000000000001',
        '00000000-0000-4000-8000-000000000001',
        'teacher.demo',
        'Demo Teacher',
        '$2y$10$x8VYzrNrLuphvD11bFS9m..bShTGlEaJAhn0DAt5hiwMPKRi60/FO',
        'school_admin',
        true,
        '2026-07-13T00:00:00Z',
        '2026-07-13T00:00:00Z'
    )
ON CONFLICT (id) DO NOTHING;
