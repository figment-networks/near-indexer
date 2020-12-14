-- +goose Up
CREATE TABLE IF NOT EXISTS epochs (
  id                 TEXT NOT NULL PRIMARY KEY,
  start_height       INTEGER NOT NULL,
  start_time         TIMESTAMP WITH TIME ZONE NOT NULL,
  end_height         INTEGER NOT NULL,
  end_time           TIMESTAMP WITH TIME ZONE NOT NULL,
  blocks_count       INTEGER DEFAULT 0,
  validators_count   INTEGER DEFAULT 0,
  average_efficiency NUMERIC
);

CREATE INDEX idx_epochs_start_time
  ON epochs(start_time);

CREATE INDEX idx_epochs_start_height
  ON epochs(start_height);

-- +goose Down
DROP TABLE epochs IF EXISTS;