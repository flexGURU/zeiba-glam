-- name: CreateOrder :one
INSERT INTO orders (user_name, user_phone_number, total_amount, status, shipping_address, payment_status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders
WHERE 
    deleted_at IS NULL
    AND (
        sqlc.narg('status')::text IS NULL 
        OR status = sqlc.narg('status')
    )
    AND (
        sqlc.narg('payment_status')::boolean IS NULL 
        OR payment_status = sqlc.narg('payment_status')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListOrdersCount :one
SELECT COUNT(*) FROM orders
WHERE 
    deleted_at IS NULL
    AND (
        sqlc.narg('status')::text IS NULL 
        OR status = sqlc.narg('status')
    )
    AND (
        sqlc.narg('payment_status')::boolean IS NULL 
        OR payment_status = sqlc.narg('payment_status')
    );

-- name: UpdateOrder :one
UPDATE orders
SET status = coalesce(sqlc.narg('status'), status),
    payment_status = coalesce(sqlc.narg('payment_status'), payment_status)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOrder :exec
UPDATE orders
SET deleted_at = now()
WHERE id = sqlc.arg('id');