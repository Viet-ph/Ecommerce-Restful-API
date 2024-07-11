-- name: CreateCartForUser :one
INSERT INTO carts (id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *; 

-- name: AddNewItemToCart :one
INSERT INTO cart_items (cart_id, product_id, quantity, price_at_time, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *; 

-- name: GetAllItemsInCart :many
SELECT sqlc.embed(ci), sqlc.embed(p), ci.quantity * ci.price_at_time AS item_total_cost
FROM cart_items ci
JOIN products p ON ci.product_id = p.id
WHERE ci.cart_id = $1;

-- name: GetItemInCart :many
SELECT sqlc.embed(ci), ci.quantity * ci.price_at_time AS item_total_cost
FROM cart_items ci
WHERE ci.cart_id = $1 AND ci.product_id = $2;

-- name: UpdateCartItem :exec
UPDATE cart_items
    SET quantity = $3,
        updated_at = NOW(),
        price_at_time = $4 
WHERE cart_id = $1 AND product_id = $2;

-- name: GetCartValueOfAllUsers :many
SELECT c.user_id, SUM(ci.quantity * ci.price_at_time) AS total_cart_value
FROM carts c
JOIN cart_items ci ON c.id = ci.cart_id
GROUP BY c.user_id
ORDER BY total_cart_value;

-- name: DeleteAllItemsInCart :exec
DELETE FROM cart_items
WHERE cart_id = (
    SELECT cart_id
    FROM carts
    WHERE user_id = $1
);

-- name: DeleteItemInCart :exec
DELETE FROM cart_items
WHERE cart_id = $1 AND product_id = $2;













