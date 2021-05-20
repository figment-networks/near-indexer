-- +goose Up
CREATE INDEX idx_validator_epochs_last_time ON validator_epochs (last_time);

-- +goose Down
DROP INDEX IF EXISTS idx_validator_epochs_last_time;
