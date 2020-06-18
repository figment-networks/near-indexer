-- +goose Up
CREATE TABLE heights (
  id          BIGSERIAL NOT NULL PRIMARY KEY,
  height      INTEGER NOT NULL,
  status      VARCHAR NOT NULL,
  retry_count INTEGER DEFAULT 0,
  error       TEXT,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX idx_height
  ON heights(height);

CREATE INDEX idx_height_status
  ON heights(status);

-- +goose Down
DROP TABLE heights;