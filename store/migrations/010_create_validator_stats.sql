-- +goose Up
CREATE TABLE validator_stats (
  id          SERIAL NOT NULL PRIMARY KEY,
  time        TIMESTAMP WITH TIME ZONE NOT NULL,
  bucket      e_interval NOT NULL,

  total_min   INTEGER,
  total_max   INTEGER,
  total_avg   INTEGER,

  active_min  INTEGER,
  active_max  INTEGER,
  active_avg  INTEGER,

  slashed_min INTEGER,
  slashed_max INTEGER,
  slashed_avg INTEGER
);

CREATE UNIQUE INDEX idx_validator_stats_bucket
  ON validator_stats(time, bucket);

-- +goose Down
DROP TABLE validator_stats;
