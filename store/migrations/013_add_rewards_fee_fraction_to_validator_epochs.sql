-- +goose Up
ALTER TABLE validator_epochs ADD COLUMN reward_fee_fraction NUMERIC(125);

-- +goose Down
ALTER TABLE validator_epochs DROP COLUMN reward_fee_fraction;