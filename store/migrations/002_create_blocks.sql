-- +goose Up
CREATE TABLE blocks (
  id                  INTEGER NOT NULL PRIMARY KEY,
  hash                VARCHAR NOT NULL,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  producer            VARCHAR NOT NULL,
  epoch               VARCHAR NOT NULL,
  gas_price           VARCHAR,
  gas_limit           INTEGER,
  gas_used            INTEGER,
  total_supply        VARCHAR,
  chunks_count        INTEGER,
  transactions_count  INTEGER,
  approvals_count     INTEGER,
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_blocks_hash
  ON blocks(hash);

CREATE INDEX idx_blocks_time
  ON blocks(time);

CREATE INDEX idx_blocks_producer
  ON blocks(producer);

CREATE INDEX idx_blocks_epoch
  ON blocks(epoch);

-- +goose Down
DROP TABLE blocks;
