-- +goose Up
ALTER TABLE refresh_tokens 
ALTER COLUMN revoked SET DEFAULT FALSE,
ALTER COLUMN revoked SET NOT NULL;

-- +goose Down
ALTER TABLE refresh_tokens 
ALTER COLUMN revoked DROP DEFAULT,
ALTER COLUMN revoked DROP NOT NULL;