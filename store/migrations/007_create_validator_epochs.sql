-- +goose Up
CREATE TABLE validator_epochs (
  id              BIGSERIAL NOT NULL PRIMARY KEY,
  account_id      VARCHAR NOT NULL,
  epoch           VARCHAR NOT NULL,
  last_height     INTEGER NOT NULL,
  last_time       TIMESTAMP WITH TIME ZONE NOT NULL,
  expected_blocks INTEGER NOT NULL,
  produced_blocks INTEGER NOT NULL,
  efficiency      NUMERIC NOT NULL
);

CREATE INDEX idx_validator_epochs_account_id
  ON validator_epochs(account_id);

CREATE UNIQUE INDEX idx_validator_epochs_account_epoch
  ON validator_epochs(account_id, epoch);

-- +goose Down
DROP TABLE validator_aggregates;
