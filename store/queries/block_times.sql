SELECT
	MIN(height) start_height,
	MAX(height) end_height,
	MIN(time) start_time,
	MAX(time) end_time,
	COUNT(*) count,
	EXTRACT(EPOCH FROM MAX(time) - MIN(time)) AS diff,
	EXTRACT(EPOCH FROM ((MAX(time) - MIN(time)) / COUNT(*))) AS avg
FROM (
	SELECT height, time
	FROM blocks
	ORDER BY height DESC
	LIMIT ?
) t
