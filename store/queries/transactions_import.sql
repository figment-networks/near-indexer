INSERT INTO transactions (
  hash,
  block_hash,
  height,
  time,
  signer,
  signer_key,
  receiver,
  signature,
  amount,
  gas_burnt,
  success,
  actions,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (hash) DO NOTHING
