-- +goose Up
CREATE TABLE IF NOT EXISTS events (
  id           SERIAL NOT NULL PRIMARY KEY,
  scope        TEXT NOT NULL,
  action       TEXT NOT NULL,
  block_height INTEGER NOT NULL,
  block_time   TIMESTAMP WITH TIME ZONE NOT NULL,
  epoch        TEXT,
  item_id      TEXT,
  item_type    TEXT,
  metadata     JSONB,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_events_scope
  ON events(scope);

CREATE INDEX idx_events_action
  ON events(action);

CREATE INDEX idx_events_block_height
  ON events(block_height);

CREATE INDEX idx_events_epoch
  ON events(epoch);

CREATE INDEX idx_events_item
  ON events(item_id, item_type);

-- +goose Down
DROP TABLE events IF EXISTS;
