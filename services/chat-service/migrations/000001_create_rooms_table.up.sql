CREATE TABLE rooms
(
    id         UUID PRIMARY KEY,
    user_a_id  UUID      NOT NULL,
    user_b_id  UUID      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);