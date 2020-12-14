SELECT
  time,
  bucket,
  blocks_count,
  block_time_avg,
  validators_count,
  transactions_count
FROM
  block_stats
WHERE
  bucket = $1
ORDER BY
  time DESC
LIMIT $2
