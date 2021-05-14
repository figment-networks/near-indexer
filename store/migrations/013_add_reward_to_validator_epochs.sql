-- +goose Up
ALTER TABLE validator_epochs ADD COLUMN reward DECIMAL(65, 0);

-- +goose Down
ALTER TABLE validator_epochs DROP COLUMN reward;