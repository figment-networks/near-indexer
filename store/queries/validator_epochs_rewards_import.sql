INSERT INTO validator_epochs_rewards (
  account_id,
  epoch,
  distributed_height,
  distributed_time,
  reward_fee,
  reward
)
VALUES @values

ON CONFLICT (account_id, epoch) DO UPDATE
SET
  distributed_height  = excluded.distributed_height,
  distributed_time    = excluded.distributed_time,
  reward_fee          = excluded.reward_fee,
  reward              = excluded.reward