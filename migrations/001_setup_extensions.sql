-- +goose Up
CREATE EXTENSION IF NOT EXISTS "timescaledb" CASCADE;

-- +goose Down
DROP EXTENSION IF EXISTS "timescaledb" CASCADE;
