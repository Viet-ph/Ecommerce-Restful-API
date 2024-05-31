-- +goose Up
ALTER TABLE refresh_tokens ADD COLUMN revoked BOOL;

-- +goose Down
ALTER TABLE refresh_tokens DROP revoked;