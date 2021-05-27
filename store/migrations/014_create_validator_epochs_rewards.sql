-- +goose Up
CREATE TABLE validator_epochs_rewards (
  id                     SERIAL NOT NULL PRIMARY KEY,
  account_id             VARCHAR NOT NULL,
  epoch                  VARCHAR NOT NULL,
  distributed_at_height  INTEGER NOT NULL,
  distributed_at_time    TIMESTAMP WITH TIME ZONE NOT NULL,
  reward_fee             INTEGER NOT NULL,
  reward                 DECIMAL(65, 0)  NOT NULL
);

CREATE UNIQUE INDEX validator_epochs_rewards_pk
  ON validator_epochs_rewards(account_id, epoch);

CREATE INDEX idx_validator_epochs_rewards_distributed_at_time
  ON validator_epochs_rewards (distributed_at_time);

-- +goose Down
DROP TABLE validator_epochs_rewards;