CREATE TABLE refresh_tokens
(
    token      TEXT PRIMARY KEY,
    user_id    BIGINT    NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
