-- name: GetDashboardStats :one
SELECT
  COUNT(*) AS total_products,
  COUNT(*) FILTER (WHERE stock_quantity > 10) AS in_stock,
  COUNT(*) FILTER (WHERE stock_quantity > 0 AND stock_quantity <= 10) AS low_stock,
  COUNT(*) FILTER (WHERE stock_quantity = 0) AS out_of_stock
FROM products
WHERE deleted_at IS NULL;