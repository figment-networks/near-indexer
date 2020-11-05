-- +goose Up
CREATE TABLE blocks (
  id                  BIGSERIAL NOT NULL PRIMARY KEY,
  height              INTEGER NOT NULL,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  hash                VARCHAR NOT NULL,
  prev_hash           VARCHAR NOT NULL,
  producer            VARCHAR NOT NULL,
  epoch               VARCHAR NOT NULL,
  gas_price           VARCHAR,
  gas_limit           INTEGER,
  gas_used            INTEGER,
  rent_paid           VARCHAR,
  validator_reward    VARCHAR,
  total_supply        VARCHAR,
  signature           VARCHAR,
  chunks_count        INTEGER,
  transactions_count  INTEGER,
  approvals_count     INTEGER,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_blocks_hash
  ON blocks(hash);

CREATE INDEX idx_blocks_height
  ON blocks(height);

CREATE INDEX idx_blocks_time
  ON blocks(time);

CREATE INDEX idx_blocks_producer
  ON blocks(producer);

CREATE INDEX idx_blocks_epoch
  ON blocks(epoch);

-- +goose Down
DROP TABLE blocks;
