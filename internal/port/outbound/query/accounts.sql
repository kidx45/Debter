-- name: GetAccountsByUserId :many
SELECT * FROM accounts
WHERE user_id = $1 ORDER BY account_type ASC;