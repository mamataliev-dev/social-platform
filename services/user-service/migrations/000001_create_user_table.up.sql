CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    nickname      VARCHAR(32)  NOT NULL UNIQUE,
    username      VARCHAR(64)  NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(100) NOT NULL,
    bio           TEXT,
    avatar_url    TEXT,
    last_login    TIMESTAMP    DEFAULT NOW(),
    created_at    TIMESTAMP    DEFAULT NOW(),
    updated_at    TIMESTAMP    DEFAULT NOW()
);
