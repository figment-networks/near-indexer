INSERT INTO delegator_epochs (
  account_id,
  validator_id,
  epoch,
  last_height,
  last_time,
  staked_balance,
  unstaked_balance,
  reward
)
VALUES @values

ON CONFLICT (account_id, validator_id, epoch) DO UPDATE
SET
  last_height         = excluded.last_height,
  last_time           = excluded.last_time,
  staked_balance      = excluded.staked_balance,
  unstaked_balance    = excluded.unstaked_balance,
  reward              = excluded.reward