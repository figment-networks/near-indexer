SELECT
	to_char(last_time, 'YYYY-MM') AS month,
	SUM(reward) AS amount
FROM
	validator_epochs
WHERE
	account_id = ?
	AND last_time BETWEEN ? AND ?
	AND reward IS NOT NULL
GROUP BY
	to_char(last_time, 'YYYY-MM')
