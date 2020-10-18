-- +goose Up
CREATE TABLE validator_counts (
  height        INTEGER NOT NULL PRIMARY KEY,
  time          TIMESTAMP WITH TIME ZONE NOT NULL,
  total_count   INTEGER,
  active_count  INTEGER,
  slashed_count INTEGER
);

CREATE INDEX idx_validator_counts_time
  ON validator_counts(time);

-- +goose Down
DROP TABLE validator_counts;
