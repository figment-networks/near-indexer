INSERT INTO blocks (
  height,
  time,
  hash,
  prev_hash,
  producer,
  epoch,
  gas_price,
  gas_limit,
  gas_used,
  rent_paid,
  validator_reward,
  total_supply,
  signature,
  chunks_count,
  transactions_count,
  approvals_count,
  created_at
)
VALUES @values
