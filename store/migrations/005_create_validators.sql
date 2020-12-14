-- +goose Up
CREATE TABLE validators (
  id              BIGSERIAL NOT NULL PRIMARY KEY,
  height          INTEGER NOT NULL,
  time            TIMESTAMP WITH TIME ZONE NOT NULL,
  account_id      VARCHAR NOT NULL,
  epoch           VARCHAR NOT NULL,
  expected_blocks INTEGER NOT NULL,
  produced_blocks INTEGER NOT NULL,
  slashed         BOOLEAN,
  stake           VARCHAR,
  efficiency      NUMERIC,
  reward_fee      INTEGER,
  created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_validators_account_id
  ON validators(account_id);

CREATE INDEX idx_validators_height
  ON validators(height);

CREATE iNDEX idx_validators_time
  ON validators(time DESC);

CREATE INDEX idx_validators_epoch
  ON validators(epoch);

-- +goose Down
DROP TABLE validators;
