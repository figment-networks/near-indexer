WITH epoch_stats AS (
  SELECT
    epoch,
    MIN(time) start_time,
    MIN(height) start_height,
    MAX(time) end_time,
    MAX(height) end_height,
    COUNT(1) blocks_count,
    COUNT(DISTINCT producer) validators_count
  FROM
    blocks
  WHERE
    epoch IN ($1)
  GROUP BY
    epoch
)
UPDATE epochs
SET
  start_time         = epoch_stats.start_time,
  start_height       = epoch_stats.start_height,
  end_time           = epoch_stats.end_time,
  end_height         = epoch_stats.end_height,
  blocks_count       = epoch_stats.blocks_count,
  validators_count   = epoch_stats.validators_count,
  average_efficiency = (
    SELECT ROUND(COALESCE(AVG(efficiency), 0), 4) FROM validator_epochs WHERE epoch IN ($1)
  )
FROM
  epoch_stats
WHERE
  epochs.uuid = epoch_stats.epoch
