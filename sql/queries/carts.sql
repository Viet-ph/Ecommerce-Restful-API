-- name: CreateCartForUser :one
INSERT INTO carts (id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *; 

-- name: AddNewItemToCart :one
INSERT INTO cart_items (id, cart_id, product_id, quantity, price_at_time, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *; 

-- name: GetAllItemsInCart :many
SELECT sqlc.embed(ci), sqlc.embed(p)
FROM cart_items ci
JOIN products p ON ci.product_id = p.id
WHERE ci.cart_id = $1;

-- name: UpdateItemQuantity :exec
UPDATE cart_items
SET quantity = $1 
WHERE id = $2;

-- name: GetCartValueByUserId :many
SELECT c.user_id, SUM(ci.quantity * ci.price_at_time) AS total_cart_value
FROM carts c
JOIN cart_items ci ON c.id = ci.cart_id
WHERE c.user_id = $1
GROUP BY c.user_id;









