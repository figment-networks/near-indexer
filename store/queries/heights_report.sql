SELECT status, COUNT(1) AS num
FROM heights
WHERE status != ''
GROUP BY status
ORDER BY num DESC
