-- +goose Up
ALTER TABLE transactions ADD COLUMN fee VARCHAR;
ALTER TABLE transactions ADD COLUMN signature VARCHAR;
ALTER TABLE transactions ADD COLUMN public_key VARCHAR;
ALTER TABLE transactions ADD COLUMN outcome JSONB;
ALTER TABLE transactions ADD COLUMN receipt JSONB;

-- +goose Down
ALTER TABLE transactions DROP COLUMN fee;
ALTER TABLE transactions DROP COLUMN signature;
ALTER TABLE transactions DROP COLUMN public_key;
ALTER TABLE transactions DROP COLUMN outcome;
ALTER TABLE transactions DROP COLUMN receipt;
