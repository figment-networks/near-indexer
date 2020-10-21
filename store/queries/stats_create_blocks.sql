INSERT INTO block_stats (
	time,
	bucket,
	blocks_count,
	block_time_avg
)
SELECT
	DATE_TRUNC('@bucket', time) AS time,
	'@bucket' AS bucket,
	COUNT(1) AS blocks_count,
	ROUND(EXTRACT(EPOCH FROM (MAX(time) - MIN(time)) / COUNT(1))::NUMERIC, 2) AS block_time_avg
FROM
	blocks
WHERE
	time >= $1 AND time <= $2
GROUP BY
	DATE_TRUNC('@bucket', time)

ON CONFLICT (time, bucket) DO UPDATE
SET
	blocks_count   = excluded.blocks_count,
	block_time_avg = excluded.block_time_avg
