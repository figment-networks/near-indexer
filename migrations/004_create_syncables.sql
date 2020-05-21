-- +goose Up
CREATE TABLE syncables (
  id           BIGSERIAL NOT NULL PRIMARY KEY,
  run_id       INTEGER NOT NULL,
  height       INTEGER NOT NULL,
  time         TIMESTAMP WITH TIME ZONE NOT NULL,
  type         VARCHAR,
  data         JSONB,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  processed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_syncables_height
  ON syncables(height, time DESC);

CREATE INDEX idx_syncables_type
  ON syncables(type, time DESC);

CREATE INDEX idx_syncables_processed_at
  ON syncables (processed_at, time DESC);

-- +goose Down
DROP TABLE syncables;