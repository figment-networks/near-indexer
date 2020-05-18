package store

var (
	syncablesRecentHeight = `
		-- DEBUG: MAX should be MIN
		--
		SELECT MAX(height) AS height
		FROM (
			SELECT DISTINCT ON(type) type, height
			FROM syncables
			WHERE processed_at IS NOT NULL
			GROUP by height, type
			ORDER BY type asc, height desc
		) t`
)
