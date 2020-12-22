-- +goose Up
CREATE TABLE validator_epochs (
  id              SERIAL NOT NULL PRIMARY KEY,
  account_id      VARCHAR NOT NULL,
  epoch           VARCHAR NOT NULL,
  last_height     INTEGER NOT NULL,
  last_time       TIMESTAMP WITH TIME ZONE NOT NULL,
  expected_blocks INTEGER NOT NULL,
  produced_blocks INTEGER NOT NULL,
  efficiency      NUMERIC NOT NULL,
  staking_balance DECIMAL(65, 0) NOT NULL
);

CREATE INDEX idx_validator_epochs_account_id
  ON validator_epochs(account_id);

CREATE UNIQUE INDEX idx_validator_epochs_account_epoch
  ON validator_epochs(account_id, epoch);

CREATE INDEX idx_validator_epochs_last_height
  ON validator_epochs(account_id, last_height);

-- +goose Down
DROP TABLE validator_epochs;
