-- +goose Up
CREATE TABLE IF NOT EXISTS epochs (
  id                 SERIAL NOT NULL PRIMARY KEY,
  uuid               TEXT NOT NULL,
  start_height       INTEGER NOT NULL,
  start_time         TIMESTAMP WITH TIME ZONE NOT NULL,
  end_height         INTEGER NOT NULL,
  end_time           TIMESTAMP WITH TIME ZONE NOT NULL,
  blocks_count       INTEGER DEFAULT 0,
  validators_count   INTEGER DEFAULT 0,
  average_efficiency NUMERIC
);

CREATE UNIQUE INDEX idx_epochs_uuid
  ON epochs(uuid);

-- +goose Down
DROP TABLE epochs IF EXISTS;
