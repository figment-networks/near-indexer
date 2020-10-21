DELETE FROM
	validator_counts
WHERE
	time <= (
		SELECT DATE_TRUNC('d', MAX(time))::timestamp - interval '48' hour
		FROM validator_counts
	)
