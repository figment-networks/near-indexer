WITH selected_blocks AS (
  SELECT id AS height, time
	FROM blocks
	ORDER BY id DESC
	LIMIT ?
)
SELECT
	MIN(height) AS start_height,
	MAX(height) AS end_height,
	MIN(time) AS start_time,
	MAX(time) AS end_time,
	COUNT(*) AS count,
	EXTRACT(EPOCH FROM MAX(time) - MIN(time)) AS diff,
	EXTRACT(EPOCH FROM ((MAX(time) - MIN(time)) / COUNT(*))) AS avg
FROM
  selected_blocks
