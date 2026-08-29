CREATE TABLE user_session (
                         id UUID primary key NOT null default gen_random_uuid(),
                         user_id uuid NOT NULL,
                         token STRING NOT NULL,
                         ttl INT DEFAULT 0,
                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);