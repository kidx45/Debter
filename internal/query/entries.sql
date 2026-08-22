-- name: GetEntriesByAccountId :many
SELECT e.* FROM entries e
JOIN accounts a ON a.id = e.account_id
WHERE e.account_id = $1 AND a.user_id = $2
ORDER BY e.created_at DESC;

-- name: FilterEntriesByDate :many
SELECT e.* FROM entries e
JOIN accounts a ON a.id = e.account_id
WHERE e.account_id = $1 AND a.user_id = $2
AND e.created_at >= $3 AND e.created_at <= $4
ORDER BY e.created_at DESC;

-- name: GetEntriesByCategoryAndType :many
SELECT e.category, SUM(e.amount) AS total
FROM entries e
JOIN accounts a ON a.id = e.account_id
WHERE e.account_id = $1 AND a.user_id = $2
AND e.type = $3
GROUP BY e.category ORDER BY e.category DESC;
