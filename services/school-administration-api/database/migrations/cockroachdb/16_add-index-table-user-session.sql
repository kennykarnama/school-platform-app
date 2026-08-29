CREATE INDEX ON user_session (token) STORING (user_id, ttl, created_at);
