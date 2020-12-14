INSERT INTO epochs (
  id,
  start_height,
  start_time,
  end_height,
  end_time,
  blocks_count,
  validators_count,
  average_efficiency
)
VALUES @values

ON CONFLICT (id) DO UPDATE
SET
  end_height         = excluded.end_height,
  end_time           = excluded.end_time,
  blocks_count       = excluded.blocks_count,
  validators_count   = excluded.validators_count,
  average_efficiency = excluded.average_efficiency
