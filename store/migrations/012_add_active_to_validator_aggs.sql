-- +goose Up
ALTER TABLE validator_aggregates ADD COLUMN active BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE validator_aggregates DROP COLUMN active;
