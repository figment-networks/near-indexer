INSERT INTO epochs (
  uuid,
  start_height,
  start_time,
  end_height,
  end_time
)
VALUES @values

ON CONFLICT (uuid) DO UPDATE
SET
  end_height = excluded.end_height,
  end_time   = excluded.end_time
