INSERT INTO transactions (
  hash,
  block_hash,
  height,
  time,
  sender,
  receiver,
  gas_burnt,
  success,
  actions,
  actions_count,
  fee,
  signature,
  public_key,
  outcome,
  receipt,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (hash) DO UPDATE
SET
  updated_at = excluded.updated_at

