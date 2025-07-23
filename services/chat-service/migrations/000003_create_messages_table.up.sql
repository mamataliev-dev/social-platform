CREATE TABLE messages
(
    id         UUID PRIMARY KEY,
    room_id    UUID      NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    sender_id  BIGINT    NOT NULL,
    content    TEXT,
    is_deleted BOOLEAN   NOT NULL DEFAULT FALSE,
    edited_at  TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);