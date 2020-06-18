-- +goose Up
CREATE TABLE validator_aggregates (
  id              BIGSERIAL NOT NULL PRIMARY KEY,
  start_height    INTEGER NOT NULL,
  start_time      TIMESTAMP WITH TIME ZONE NOT NULL,
  last_height     INTEGER NOT NULL,
  last_time       TIMESTAMP WITH TIME ZONE NOT NULL,
  account_id      VARCHAR NOT NULL,
  expected_blocks INTEGER NOT NULL,
  produced_blocks INTEGER NOT NULL,
  slashed         BOOLEAN,
  stake           VARCHAR,
  efficiency      NUMERIC,
  created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_account_id
  ON validator_aggregates(account_id);

-- +goose Down
DROP TABLE validator_aggregates;