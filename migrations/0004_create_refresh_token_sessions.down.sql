DROP INDEX IF EXISTS idx_refresh_tokens_session_id;
DROP INDEX IF EXISTS idx_refresh_token_sessions_user_id;

ALTER TABLE refresh_token_sessions
    DROP CONSTRAINT IF EXISTS fk_refresh_token_sessions_current_token_id;

ALTER TABLE refresh_tokens
    DROP CONSTRAINT IF EXISTS fk_refresh_tokens_session_id;

ALTER TABLE refresh_tokens
    DROP COLUMN IF EXISTS session_id;

DROP TABLE IF EXISTS refresh_token_sessions;
