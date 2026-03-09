-- +goose Up
CREATE TABLE outbox
(
    message_id UUID PRIMARY KEY,
    message    JSONB
);

-- +goose Down
DROP TABLE IF EXISTS outbox