INSERT INTO validator_stats (
  time,
  bucket,
  total_min, total_max, total_avg,
  active_min, active_max, active_avg,
  slashed_min, slashed_max, slashed_avg
)
SELECT
  DATE_TRUNC('@bucket', time) AS time,
  '@bucket' AS bucket,
  MIN(total_count),   MAX(total_count),   ROUND(AVG(total_count), 2),
  MIN(active_count),  MAX(active_count),  ROUND(AVG(active_count), 2),
  MIN(slashed_count), MAX(slashed_count), ROUND(AVG(slashed_count), 2)
FROM
  validator_counts
WHERE
  time >= $1 AND time <= $2
GROUP BY
  DATE_TRUNC('@bucket', time)

ON CONFLICT (time, bucket) DO UPDATE
SET
  total_min   = excluded.total_min,
  total_max   = excluded.total_max,
  total_avg   = excluded.total_avg,
  active_min  = excluded.active_min,
  active_max  = excluded.active_max,
  active_avg  = excluded.active_avg,
  slashed_min = excluded.slashed_min,
  slashed_max = excluded.slashed_max,
  slashed_avg = excluded.slashed_avg
