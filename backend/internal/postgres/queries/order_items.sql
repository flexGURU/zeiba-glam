-- name: ListOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;