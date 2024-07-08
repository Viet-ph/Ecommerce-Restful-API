-- name: CreateProduct :one
INSERT INTO products (id, 
                    title, 
                    supplier, 
                    category, 
                    price, 
                    image_url, 
                    description, 
                    product_location, 
                    created_at, 
                    updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: ListProducts :many
SELECT *
FROM products
WHERE 
    (sqlc.narg(supplier)::TEXT IS NULL OR supplier = sqlc.narg(supplier)) AND
    (sqlc.narg(category)::TEXT IS NULL OR category = sqlc.narg(category)) AND
    (sqlc.narg(product_location)::TEXT IS NULL OR product_location = sqlc.narg(product_location));

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: DeleteProductById :exec
DELETE FROM products WHERE id = $1;