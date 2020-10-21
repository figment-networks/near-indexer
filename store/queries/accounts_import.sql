INSERT INTO accounts (
  name,
  start_height,
  start_time,
  last_height,
  last_time,
  balance,
  staking_balance,
  created_at,
  updated_at
)
VALUES @values

ON CONFLICT (name) DO UPDATE
SET
  last_height     = excluded.last_height,
  last_time       = excluded.last_time,
  balance         = excluded.balance,
  staking_balance = excluded.staking_balance,
  updated_at      = excluded.updated_at
