-- +goose Up
ALTER TABLE validator_epochs ADD COLUMN reward_fee INTEGER;

-- +goose Down
ALTER TABLE validator_epochs DROP COLUMN reward_fee;
