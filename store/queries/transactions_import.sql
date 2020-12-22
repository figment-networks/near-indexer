INSERT INTO transactions (
  hash,
  block_hash,
  height,
  time,
  sender,
  receiver,
  amount,
  gas_burnt,
  success,
  actions,
  actions_count,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (hash) DO UPDATE
SET
  updated_at = excluded.updated_at

