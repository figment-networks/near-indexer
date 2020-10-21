SELECT
	validator_aggregates.*,
	{{ array }}
		SELECT
			validator_epochs.epoch,
			validator_epochs.last_height,
			validator_epochs.last_time,
			validator_epochs.expected_blocks,
			validator_epochs.produced_blocks,
			validator_epochs.efficiency
	{{ end_array }} AS epochs
FROM
	validator_aggregates
LEFT JOIN validator_epochs
	ON validator_epochs.account_id = validator_aggregates.account_id
WHERE
	validator_aggregates.account_id = ?
LIMIT 1
