-- +goose Up
CREATE TABLE IF NOT EXISTS jwt_session (
    token_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    access_expires_at TIMESTAMPTZ NOT NULL,
    refresh_expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS jwt_session_user_id_idx
    ON jwt_session (user_id);

CREATE INDEX IF NOT EXISTS jwt_session_refresh_expires_at_idx
    ON jwt_session (refresh_expires_at);

-- +goose Down
DROP TABLE IF EXISTS jwt_session;
