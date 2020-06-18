-- +goose Up
CREATE TABLE accounts (
  id              BIGSERIAL NOT NULL PRIMARY KEY,
  name            VARCHAR NOT NULL,
  start_height    INTEGER NOT NULL,
  start_time      TIMESTAMP WITH TIME ZONE NOT NULL,
  last_height     INTEGER NOT NULL,
  last_time       TIMESTAMP WITH TIME ZONE NOT NULL,
  balance         VARCHAR,
  staking_balance VARCHAR,
  created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_accounts_name
  ON accounts(name);

-- +goose Down
DROP TABLE accounts;