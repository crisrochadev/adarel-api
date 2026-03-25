CREATE TABLE IF NOT EXISTS refresh_token_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_token_id BIGINT,
    device_fingerprint TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE refresh_tokens
    ADD COLUMN IF NOT EXISTS session_id BIGINT;

ALTER TABLE refresh_tokens
    ADD CONSTRAINT fk_refresh_tokens_session_id
    FOREIGN KEY (session_id) REFERENCES refresh_token_sessions(id) ON DELETE CASCADE;

ALTER TABLE refresh_token_sessions
    ADD CONSTRAINT fk_refresh_token_sessions_current_token_id
    FOREIGN KEY (current_token_id) REFERENCES refresh_tokens(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_refresh_token_sessions_user_id ON refresh_token_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_session_id ON refresh_tokens(session_id);
