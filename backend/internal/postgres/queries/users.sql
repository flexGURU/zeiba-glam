-- name: CreateUser :one
INSERT INTO users (name, email, phone_number, refresh_token, password, is_admin)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT 
    id,
    name,
    email,
    phone_number,
    is_admin,
    created_at
FROM users
WHERE
    (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search') 
        OR LOWER(email) LIKE sqlc.narg('search') 
        OR LOWER(phone_number) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('is_admin')::boolean IS NULL 
        OR is_admin = sqlc.narg('is_admin')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListUsersCount :one
SELECT COUNT(*) AS total_users
FROM users
WHERE
    (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search') 
        OR LOWER(email) LIKE sqlc.narg('search') 
        OR LOWER(phone_number) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('is_admin')::boolean IS NULL 
        OR is_admin = sqlc.narg('is_admin')
    );

-- name: UpdateUser :one
UPDATE users 
SET name = coalesce(sqlc.narg('name'), name),
    email = coalesce(sqlc.narg('email'), email),
    phone_number = coalesce(sqlc.narg('phone_number'), phone_number),
    password = coalesce(sqlc.narg('password'), password),
    is_admin = coalesce(sqlc.narg('is_admin'), is_admin),
    refresh_token = coalesce(sqlc.narg('refresh_token'), refresh_token)
WHERE id = sqlc.arg('id')
RETURNING *;

