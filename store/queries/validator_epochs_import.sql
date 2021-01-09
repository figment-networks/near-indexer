INSERT INTO validator_epochs (
  account_id,
  epoch,
  last_height,
  last_time,
  expected_blocks,
  produced_blocks,
  efficiency,
  staking_balance,
  reward_fee
)
VALUES @values

ON CONFLICT (account_id, epoch) DO UPDATE
SET
  last_height     = excluded.last_height,
  last_time       = excluded.last_time,
  expected_blocks = excluded.expected_blocks,
  produced_blocks = excluded.produced_blocks,
  efficiency      = ROUND(excluded.efficiency, 4),
  staking_balance = excluded.staking_balance,
  reward_fee      = COALESCE(excluded.reward_fee, validator_epochs.reward_fee)
