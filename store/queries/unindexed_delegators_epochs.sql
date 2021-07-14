SELECT
  id,
  start_height,
  start_time,
  end_height,
  end_time,
  blocks_count,
  validators_count,
  average_efficiency
FROM epochs
WHERE start_height < (
        SELECT start_height FROM epochs
         INNER JOIN delegator_epochs ON epochs.id = delegator_epochs.epoch
         ORDER BY start_height ASC LIMIT 1
    )
ORDER BY start_height DESC
