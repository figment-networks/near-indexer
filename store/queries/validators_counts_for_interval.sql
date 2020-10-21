SELECT
	time AS time_interval,
	active_avg AS count
FROM
	validator_stats
WHERE
	bucket = $1
ORDER BY
	time DESC
LIMIT $2
