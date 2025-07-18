-- name: CreateProduct :one
INSERT INTO products (name, description, price, category, image_url, size, color, stock_quantity, updated_by, sub_category)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductsByCategory :many
SELECT * FROM products WHERE category = $1;

-- name: GetProductsBySubCategory :many
SELECT * FROM products WHERE sub_category = $1;

-- name: UpdateProductCategory :exec
UPDATE products
SET category = sqlc.arg('new_category')
WHERE category = sqlc.arg('old_category');

-- name: UpdateProductSubCategory :exec
UPDATE products
SET sub_category = sqlc.arg('new_sub_category')
WHERE sub_category = sqlc.arg('old_sub_category');

-- name: ListProducts :many
SELECT * FROM products
WHERE
    deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('price_from')::float IS NULL OR price >= sqlc.narg('price_from')
    )
    AND (
        sqlc.narg('price_to')::float IS NULL OR price <= sqlc.narg('price_to')
    )
    AND (
        sqlc.narg('category')::text[] IS NULL OR category = ANY(sqlc.narg('category')::text[])
    )
    AND (
        sqlc.narg('sub_category')::text[] IS NULL OR sub_category = ANY(sqlc.narg('sub_category')::text[])
    )
    AND (
        sqlc.narg('size')::text[] IS NULL OR size && sqlc.narg('size')
    )
    AND (
        sqlc.narg('color')::text[] IS NULL OR color && sqlc.narg('color')
    )
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListProductsCount :one
SELECT COUNT(*) AS total_products
FROM products
WHERE
    deleted_at IS NULL
    AND (
        COALESCE(sqlc.narg('search'), '') = '' 
        OR LOWER(name) LIKE sqlc.narg('search')
    )
    AND (
        sqlc.narg('price_from')::float IS NULL OR price >= sqlc.narg('price_from')
    )
    AND (
        sqlc.narg('price_to')::float IS NULL OR price <= sqlc.narg('price_to')
    )
    AND (
        sqlc.narg('category')::text[] IS NULL OR category = ANY(sqlc.narg('category')::text[])
    )
    AND (
        sqlc.narg('sub_category')::text[] IS NULL OR sub_category = ANY(sqlc.narg('sub_category')::text[])
    )
    AND (
        sqlc.narg('size')::text[] IS NULL OR size && sqlc.narg('size')
    )
    AND (
        sqlc.narg('color')::text[] IS NULL OR color && sqlc.narg('color')
    ); 
    
-- name: UpdateProduct :one
UPDATE products
SET updated_by = sqlc.arg('updated_by'),
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    price = coalesce(sqlc.narg('price'), price),
    category = coalesce(sqlc.narg('category'), category),
    sub_category = coalesce(sqlc.narg('sub_category'), sub_category),
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    size = coalesce(sqlc.narg('size'), size),
    color = coalesce(sqlc.narg('color'), color),
    stock_quantity = coalesce(sqlc.narg('stock_quantity'), stock_quantity)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = NOW()
WHERE id = $1;