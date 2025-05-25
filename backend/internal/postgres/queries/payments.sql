-- name: CreatePayment :one
INSERT INTO payments (order_id, amount, transaction_id, payment_method, payment_status, paid_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPaymentsOverviewByOrderID :many
SELECT * FROM payments WHERE order_id = $1;

-- name: GetPayment :one
SELECT 
    p.*,
    o.id as "order_id",
    o.user_name as "order_user_name",
    o.user_phone_number as "order_user_phone_number", 
    o.total_amount as "order_total_amount",
    o.status as "order_status",
    o.shipping_address as "order_shipping_address",
    o.payment_status as "order_payment_status",
    o.created_at as "order_created_at"
FROM payments p
JOIN orders o ON o.id = p.order_id 
WHERE 
    (
      sqlc.narg('id') IS NOT NULL AND id = sqlc.narg('id')
      OR
      sqlc.narg('order_id') IS NOT NULL AND order_id = sqlc.narg('order_id')
    )
LIMIT 1;

-- -- name: GetPaymentByOrderID :one
-- SELECT 
--     p.*,
--     o.id as "order.id",
--     o.user_name as "order.user_name",
--     o.user_phone_number as "order.user_phone_number", 
--     o.total_amount as "order.total_amount",
--     o.status as "order.status",
--     o.shipping_address as "order.shipping_address",
--     o.payment_status as "order.payment_status",
--     o.created_at as "order.created_at"
-- FROM payments p
-- JOIN orders o ON o.id = p.order_id 
-- WHERE p.order_id = $1;

-- name: ListPayments :many
SELECT *
FROM payments
WHERE 
    (
        sqlc.narg('payment_method')::text IS NULL 
        OR payment_method = sqlc.narg('payment_method')
    )
    AND (
        sqlc.narg('payment_status')::boolean IS NULL 
        OR payment_status = sqlc.narg('payment_status')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListPaymentsCount :one
SELECT COUNT(*) AS total_payments
FROM payments 
WHERE 
    (
        sqlc.narg('payment_method')::text IS NULL 
        OR payment_method = sqlc.narg('payment_method')
    )
    AND (
        sqlc.narg('payment_status')::boolean IS NULL 
        OR payment_status = sqlc.narg('payment_status')
    );

-- name: UpdatePayment :one
UPDATE payments 
SET payment_status = coalesce(sqlc.narg('payment_status'), payment_status), 
    paid_at = coalesce(sqlc.narg('paid_at'), paid_at)
WHERE id = sqlc.arg('id')
RETURNING *;