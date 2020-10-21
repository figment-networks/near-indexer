INSERT INTO validator_epochs (
	account_id,
	epoch,
	last_height,
	last_time,
	expected_blocks,
	produced_blocks,
	efficiency
)
VALUES @values

ON CONFLICT (account_id, epoch) DO UPDATE
SET
	last_height     = excluded.last_height,
	last_time       = excluded.last_time,
	expected_blocks = excluded.expected_blocks,
	produced_blocks = excluded.produced_blocks,
	efficiency      = excluded.efficiency;
