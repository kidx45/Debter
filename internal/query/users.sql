-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users 
WHERE username = $1;

-- name: UpdateFullNameByUsername :one
UPDATE users SET full_name = $1
WHERE username = $2
RETURNING *;

-- name: UpdateUserNameByUsername :one
UPDATE users SET username = sqlc.arg(new_username)
WHERE username = $1
RETURNING *;

-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE username = $1;