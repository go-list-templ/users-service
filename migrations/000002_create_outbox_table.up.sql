CREATE TABLE IF NOT EXISTS outbox
(
    message_id UUID PRIMARY KEY,
    message    JSONB
);