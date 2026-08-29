-- Local development login: teacher.demo / school123
INSERT INTO teacher (id, alternative_id, name, password, created_at, updated_at) VALUES
    (
        '30000000-0000-4000-8000-000000000001',
        'teacher.demo',
        'Demo Teacher',
        '$2y$10$x8VYzrNrLuphvD11bFS9m..bShTGlEaJAhn0DAt5hiwMPKRi60/FO',
        '2026-07-13T00:00:00Z',
        '2026-07-13T00:00:00Z'
    )
ON CONFLICT (id) DO NOTHING;
