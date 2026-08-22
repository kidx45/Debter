-- name: CreateSession :one
INSERT INTO sessions (user_id, refresh_token, user_agent, client_ip, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
WHERE refresh_token = $1;

-- name: UpdateSessionRefreshToken :one
UPDATE sessions SET refresh_token = $2, expires_at = $3
WHERE id = $1
RETURNING *;
