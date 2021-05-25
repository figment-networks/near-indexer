INSERT INTO delegator_epochs (
  account_id,
  validator_id,
  epoch,
  distributed_at_height,
  distributed_at_time,
  staked_balance,
  unstaked_balance,
  reward
)
VALUES @values

ON CONFLICT (account_id, validator_id, epoch) DO UPDATE
SET
  distributed_at_height  = excluded.distributed_at_height,
  distributed_at_time    = excluded.distributed_at_time,
  staked_balance         = excluded.staked_balance,
  unstaked_balance       = excluded.unstaked_balance,
  reward                 = excluded.reward