SELECT
	to_char(distributed_time, $INTERVAL) AS interval,
	SUM(reward) AS amount
FROM
	validator_epochs_rewards
WHERE
	account_id = ?
	AND distributed_time BETWEEN ? AND ?
	AND reward IS NOT NULL
GROUP BY
	to_char(distributed_time, $INTERVAL)
