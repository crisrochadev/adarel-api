ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ;

ALTER TABLE refresh_tokens
    ADD CONSTRAINT chk_refresh_tokens_expiration
    CHECK (expires_at > issued_at);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_active_by_user
    ON refresh_tokens(user_id, expires_at)
    WHERE revoked_at IS NULL;
