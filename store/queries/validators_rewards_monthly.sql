SELECT
	date_trunc('month', last_time) AS month,
	SUM(reward) AS amount
FROM
	validator_epochs
WHERE
	account_id = ?
	AND last_time BETWEEN ? AND ?
	AND reward IS NOT NULL
GROUP BY
	date_trunc('month', last_time)
