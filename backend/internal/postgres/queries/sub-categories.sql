-- name: CreateSubCategory :one
INSERT INTO sub_categories (category_id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSubCategoryByID :one
SELECT sub_categories.*, 
       categories.name AS category_name, 
       categories.description AS category_description 
FROM sub_categories 
JOIN categories ON sub_categories.category_id = categories.id
WHERE sub_categories.id = $1;

-- name: GetSubCategoryByCategoryIDAndID :one
SELECT * FROM sub_categories
WHERE category_id = $1 AND id = $2;

-- name: ListSubCategoriesByCategory :many
SELECT * FROM sub_categories
WHERE category_id = $1
ORDER BY name;

-- name: ListSubCategories :many
SELECT sub_categories.*, 
       categories.name AS category_name, 
       categories.description AS category_description 
FROM sub_categories
JOIN categories ON sub_categories.category_id = categories.id
ORDER BY sub_categories.name;

-- name: UpdateSubCategory :one
UPDATE sub_categories
SET name = COALESCE(sqlc.narg('name'), name),
    category_id = COALESCE(sqlc.narg('category_id'), category_id),
    description = COALESCE(sqlc.narg('description'), description)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSubCategory :exec
DELETE FROM sub_categories
WHERE id = $1;