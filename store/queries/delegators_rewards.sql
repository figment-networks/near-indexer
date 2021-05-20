SELECT
	to_char(distributed_time, $INTERVAL) AS interval,
	SUM(reward) AS amount
FROM
	delegator_epochs
WHERE
	account_id = ?
	AND validator_id = ?
	AND distributed_time BETWEEN ? AND ?
GROUP BY
	to_char(distributed_time, $INTERVAL)
