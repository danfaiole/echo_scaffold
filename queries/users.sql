-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  username, password, email
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
  set password = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
