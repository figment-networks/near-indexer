-- +goose Up
CREATE TYPE e_interval AS ENUM ('h', 'd', 'w');

-- +goose Down
DROP TYPE e_interval;