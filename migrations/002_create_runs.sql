-- +goose Up
CREATE TABLE runs (
  id         BIGSERIAL NOT NULL PRIMARY KEY,
  height     INTEGER NOT NULL,
  success    BOOLEAN NOT NULL,
  error      TEXT,
  duration   NUMERIC,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_runs_height ON runs(height);

-- +goose Down
DROP TABLE runs;