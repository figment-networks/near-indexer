INSERT INTO validators (
  height,
  time,
  account_id,
  epoch,
  expected_blocks,
  produced_blocks,
  slashed,
  stake,
  efficiency,
  reward_fee,
  created_at,
  updated_at
)
VALUES @values
