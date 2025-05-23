-- name: CreateProduct :one
INSERT INTO products (name, description, price, category, image_url, size, color, stock_quantity, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products WHERE deleted_at IS NULL;

-- name: UpdateProduct :one
UPDATE products
SET name = $1,
    description = $2,
    price = $3,
    category = $4,
    image_url = $5,
    size = $6,
    color = $7,
    stock_quantity = $8
WHERE id = $9
RETURNING *;

-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = NOW()
WHERE id = $1;