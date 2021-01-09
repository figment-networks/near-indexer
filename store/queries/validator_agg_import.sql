INSERT INTO validator_aggregates(
  start_height,
  start_time,
  last_height,
  last_time,
  account_id,
  expected_blocks,
  produced_blocks,
  slashed,
  stake,
  efficiency,
  active,
  reward_fee,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT(account_id) DO UPDATE
SET
  last_height     = excluded.last_height,
  last_time       = excluded.last_time,
  expected_blocks = COALESCE((SELECT SUM(expected_blocks) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),
  produced_blocks = COALESCE((SELECT SUM(produced_blocks) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),
  efficiency      = COALESCE((SELECT AVG(efficiency) FROM validator_epochs WHERE account_id = excluded.account_id LIMIT 1), 0),
  stake           = excluded.stake,
  slashed         = excluded.slashed,
  active          = excluded.active,
  reward_fee      = excluded.reward_fee,
  updated_at      = excluded.updated_at
