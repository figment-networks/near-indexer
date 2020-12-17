WITH recent_heights AS (
  SELECT DISTINCT height FROM validators ORDER BY height DESC LIMIT ?
)
DELETE FROM validators
WHERE height < (SELECT COALESCE(MIN(height), 0) FROM recent_heights)