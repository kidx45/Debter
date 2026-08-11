-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users 
WHERE username = $1;

-- name: UpdateUserNameByUsername :exec
UPDATE users SET full_name = $1
WHERE username = $2;

-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE username = $1;