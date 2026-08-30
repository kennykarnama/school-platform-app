-- +goose Up
CREATE INDEX IF NOT EXISTS user_session_token_idx
    ON user_session (token) STORING (user_id, ttl, created_at);

-- +goose Down
DROP INDEX IF EXISTS user_session_token_idx;
