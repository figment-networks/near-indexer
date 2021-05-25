INSERT INTO validator_epochs_rewards (
  account_id,
  epoch,
  distributed_at_height,
  distributed_at_time,
  reward_fee,
  reward
)
VALUES @values

ON CONFLICT (account_id, epoch) DO UPDATE
SET
  distributed_at_height  = excluded.distributed_at_height,
  distributed_at_time    = excluded.distributed_at_time,
  reward_fee             = excluded.reward_fee,
  reward                 = excluded.reward