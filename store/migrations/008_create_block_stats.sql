-- +goose Up
CREATE TYPE e_interval AS ENUM ('h', 'd', 'w');

CREATE TABLE block_stats (
  id                  SERIAL NOT NULL PRIMARY KEY,
  time                TIMESTAMP WITH TIME ZONE NOT NULL,
  bucket              e_interval NOT NULL,

  blocks_count        INTEGER,
  block_time_avg      NUMERIC,
  validators_count    INTEGER
);

CREATE UNIQUE INDEX idx_block_stats_bucket
  ON block_stats(time, bucket);

-- +goose Down
DROP TABLE block_stats;
DROP TYPE e_interval;
