CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    nickname      VARCHAR(32)  NOT NULL UNIQUE,
    user_name     VARCHAR(64)  NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(100) NOT NULL,
    bio           TEXT,
    avatar_url    TEXT,
    last_login_at TIMESTAMP             DEFAULT NOW(),
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP             DEFAULT NOW()
);
