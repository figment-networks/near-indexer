INSERT INTO delegator_epochs (
  account_id,
  validator_id,
  epoch,
  distributed_height,
  distributed_time,
  staked_balance,
  unstaked_balance,
  reward
)
VALUES @values

ON CONFLICT (account_id, validator_id, epoch) DO UPDATE
SET
  distributed_height  = excluded.distributed_height,
  distributed_time    = excluded.distributed_time,
  staked_balance      = excluded.staked_balance,
  unstaked_balance    = excluded.unstaked_balance,
  reward              = excluded.reward