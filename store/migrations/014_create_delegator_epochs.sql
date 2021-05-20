-- +goose Up
CREATE TABLE delegator_epochs (
  id                     SERIAL NOT NULL PRIMARY KEY,
  account_id             VARCHAR NOT NULL, -- delegator
  validator_id           VARCHAR NOT NULL,
  epoch                  VARCHAR NOT NULL,
  distributed_height     INTEGER NOT NULL,
  distributed_time       TIMESTAMP WITH TIME ZONE NOT NULL,
  staked_balance         DECIMAL(65, 0) NOT NULL,
  unstaked_balance       DECIMAL(65, 0) NOT NULL,
  reward                 DECIMAL(65, 0)  NOT NULL
);

CREATE INDEX idx_delegator_epochs_account_id
  ON delegator_epochs(account_id);

CREATE UNIQUE INDEX idx_delegator_epochs_account_epoch
  ON delegator_epochs(account_id, validator_id, epoch);

CREATE INDEX idx_delegator_epochs_distributed_time ON delegator_epochs (distributed_time);

-- +goose Down
DROP TABLE delegator_epochs;
