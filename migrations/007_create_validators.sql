-- +goose Up
CREATE TABLE validators (
  id              BIGSERIAL NOT NULL,
  height          INTEGER NOT NULL,
  time            TIMESTAMP WITH TIME ZONE NOT NULL,
  account_id      VARCHAR NOT NULL,
  epoch           VARCHAR NOT NULL,
  expected_blocks INTEGER NOT NULL,
  produced_blocks INTEGER NOT NULL,
  slashed         BOOLEAN,
  stake           VARCHAR,
  efficiency      NUMERIC,
  created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at      TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (time, id)
);

SELECT create_hypertable('validators', 'time', if_not_exists => TRUE);

CREATE INDEX idx_validators_account_id
  ON validators(account_id);

CREATE INDEX idx_validators_height
  ON validators(height);

CREATE iNDEX idx_validators_height_time
  ON validators(height, time DESC);

-- +goose Down
DROP TABLE validators;