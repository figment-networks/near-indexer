SELECT
    to_char(distributed_at_time, $INTERVAL) AS interval,
    SUM(reward) AS amount
FROM
    validator_epochs_rewards
WHERE account_id = ?
    AND distributed_at_time BETWEEN ? AND ?
    AND reward IS NOT NULL
GROUP BY
    to_char(distributed_at_time, $INTERVAL)
