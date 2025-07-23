CREATE TABLE rooms
(
    id             UUID PRIMARY KEY,
    initiator_id   BIGINT    NOT NULL,
    participant_id BIGINT    NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP
);