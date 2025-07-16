CREATE TABLE message_receipts
(
    message_id UUID NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    user_id    UUID NOT NULL,
    seen       BOOLEAN DEFAULT FALSE,
    seen_at    TIMESTAMP,
    PRIMARY KEY (message_id, user_id)
);
