-- +goose Up
CREATE TABLE blocks (
  id                  BIGSERIAL NOT NULL,
  height              INTEGER NOT NULL,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  hash                VARCHAR NOT NULL,
  prev_hash           VARCHAR NOT NULL,
  producer            VARCHAR NOT NULL,
  gas_price           VARCHAR,
  gas_limit           INTEGER,
  gas_used            INTEGER,
  rent_paid           VARCHAR,
  validator_reward    VARCHAR,
  total_supply        VARCHAR,
  signature           VARCHAR,
  chunks_count        INTEGER,
  transactions_count  INTEGER,
  app_version         VARCHAR,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,

  PRIMARY KEY (id, time)
);

SELECT create_hypertable('blocks', 'time', if_not_exists => TRUE);

CREATE INDEX idx_blocks_hash
  ON blocks(time, hash);

CREATE INDEX idx_blocks_height
  ON blocks(time, height);

CREATE INDEX idx_blocks_producer
  ON blocks(producer);

-- +goose Down
DROP TABLE blocks;