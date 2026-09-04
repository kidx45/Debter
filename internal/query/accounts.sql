-- name: GetAccountsByUserId :many
SELECT * FROM accounts
WHERE user_id = $1 ORDER BY account_type ASC;

-- name: CreateAccount :one
INSERT INTO accounts (user_id, account_type, account_number, balance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1;

-- name: DebitAccount :execrows
UPDATE accounts
SET balance = balance + $1
WHERE id = $2 AND user_id = $3 AND balance + $1 >= 0;

-- name: CreditAccount :exec
UPDATE accounts
SET balance = balance + $1
WHERE id = $2 AND user_id = $3;