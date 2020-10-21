INSERT INTO validator_counts (
  height,
  time,
  total_count,
  active_count,
  slashed_count
)
SELECT
  blocks.height,
  blocks.time,
  (SELECT COUNT(1) FROM validators WHERE validators.height = blocks.height) total_validators,
  (SELECT COUNT(1) FROM validators WHERE validators.height = blocks.height AND slashed = false AND efficiency > 0) active_validators,
  (SELECT COUNT(1) FROM validators WHERE validators.height = blocks.height AND slashed = true) slashed_validators
FROM
  blocks
WHERE
  blocks.height = $1

ON CONFLICT (height) DO UPDATE
SET
  time          = excluded.time,
  total_count   = excluded.total_count,
  active_count  = excluded.active_count,
  slashed_count = excluded.slashed_count;
