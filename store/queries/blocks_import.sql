INSERT INTO blocks (
  id,
  time,
  hash,
  producer,
  epoch,
  gas_price,
  gas_limit,
  gas_used,
  total_supply,
  chunks_count,
  transactions_count,
  approvals_count,
  created_at
)
VALUES @values

ON CONFLICT (id) DO NOTHING