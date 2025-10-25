-- name: CreateUser :one
INSERT INTO users (fullname, email, password)
VALUES ($1, $2, $3)
RETURNING id, fullname, email, created_at, updated_at;
-- name: GetUserByEmail :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, fullname, created_at, updated_at
FROM users
WHERE id = $1;
