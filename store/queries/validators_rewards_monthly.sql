SELECT
	date_trunc('month', last_height), 'MM/YYYY' AS month,
	SUM(rewards) AS amount
FROM
	validator_epochs
WHERE
	account_id = ?
	AND last_height BETWEEN ? AND ?
	AND reward IS NOT NULL
GROUP BY
	date_trunc('month', last_height)