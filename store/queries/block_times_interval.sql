SELECT
  time AS time_interval,
  blocks_count AS count,
  block_time_avg AS avg
FROM
  block_stats
WHERE
  bucket = $1
ORDER BY
  time DESC
LIMIT $2
