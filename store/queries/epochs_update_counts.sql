WITH epoch_stats AS (
  SELECT
    epoch,
    MIN(time) AS start_time,
    MIN(id) AS start_height,
    MAX(time) AS end_time,
    MAX(id) AS end_height,
    COUNT(1) AS blocks_count,
    COUNT(DISTINCT producer) AS validators_count
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
  epochs.id = epoch_stats.epoch
