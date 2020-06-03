-- +goose Up
CREATE TABLE validator_counts (
  height        INTEGER NOT NULL,
  time          TIMESTAMP WITH TIME ZONE NOT NULL,
  total_count   INTEGER,
  active_count  INTEGER,
  slashed_count INTEGER,

  PRIMARY KEY(height)
);

CREATE INDEX idx_validator_counts_time
  ON validator_counts(time);

-- +goose Down
DROP TABLE validator_counts;
