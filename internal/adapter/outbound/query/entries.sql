-- name: GetEntriesByAccountId :many
SELECT * FROM entries 
WHERE account_id = $1 ORDER BY created_at DESC;

-- name: FilterEntriesByDate :many
SELECT * FROM entries
WHERE account_id = $1 AND created_at >= $2 AND created_at <= $3
ORDER BY created_at DESC;

-- name: GetEntriesByCategoryAndType :many
SELECT category, SUM(amount) AS total
FROM entries
WHERE account_id = $1 AND
type = $2
GROUP BY category ORDER BY category DESC;