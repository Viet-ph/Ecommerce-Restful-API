// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: carts.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addNewItemToCart = `-- name: AddNewItemToCart :one
INSERT INTO cart_items (cart_id, product_id, quantity, price_at_time, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING cart_id, product_id, quantity, price_at_time, created_at, updated_at
`

type AddNewItemToCartParams struct {
	CartID      uuid.UUID
	ProductID   uuid.UUID
	Quantity    int32
	PriceAtTime string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) AddNewItemToCart(ctx context.Context, arg AddNewItemToCartParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, addNewItemToCart,
		arg.CartID,
		arg.ProductID,
		arg.Quantity,
		arg.PriceAtTime,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i CartItem
	err := row.Scan(
		&i.CartID,
		&i.ProductID,
		&i.Quantity,
		&i.PriceAtTime,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createCartForUser = `-- name: CreateCartForUser :one
INSERT INTO carts (id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, created_at, updated_at
`

type CreateCartForUserParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateCartForUser(ctx context.Context, arg CreateCartForUserParams) (Cart, error) {
	row := q.db.QueryRowContext(ctx, createCartForUser,
		arg.ID,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllItemsInCart = `-- name: DeleteAllItemsInCart :exec
DELETE FROM cart_items
WHERE cart_id = (
    SELECT cart_id
    FROM carts
    WHERE user_id = $1
)
`

func (q *Queries) DeleteAllItemsInCart(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteAllItemsInCart, userID)
	return err
}

const deleteItemInCart = `-- name: DeleteItemInCart :exec
DELETE FROM cart_items
WHERE cart_id = $1 AND product_id = $2
`

type DeleteItemInCartParams struct {
	CartID    uuid.UUID
	ProductID uuid.UUID
}

func (q *Queries) DeleteItemInCart(ctx context.Context, arg DeleteItemInCartParams) error {
	_, err := q.db.ExecContext(ctx, deleteItemInCart, arg.CartID, arg.ProductID)
	return err
}

const getAllItemsInCart = `-- name: GetAllItemsInCart :many
SELECT ci.cart_id, ci.product_id, ci.quantity, ci.price_at_time, ci.created_at, ci.updated_at, p.id, p.title, p.supplier, p.category, p.price, p.image_url, p.description, p.product_location, p.created_at, p.updated_at, ci.quantity * ci.price_at_time AS item_total_cost
FROM cart_items ci
JOIN products p ON ci.product_id = p.id
WHERE ci.cart_id = $1
`

type GetAllItemsInCartRow struct {
	CartItem      CartItem
	Product       Product
	ItemTotalCost int32
}

func (q *Queries) GetAllItemsInCart(ctx context.Context, cartID uuid.UUID) ([]GetAllItemsInCartRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllItemsInCart, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllItemsInCartRow
	for rows.Next() {
		var i GetAllItemsInCartRow
		if err := rows.Scan(
			&i.CartItem.CartID,
			&i.CartItem.ProductID,
			&i.CartItem.Quantity,
			&i.CartItem.PriceAtTime,
			&i.CartItem.CreatedAt,
			&i.CartItem.UpdatedAt,
			&i.Product.ID,
			&i.Product.Title,
			&i.Product.Supplier,
			&i.Product.Category,
			&i.Product.Price,
			&i.Product.ImageUrl,
			&i.Product.Description,
			&i.Product.ProductLocation,
			&i.Product.CreatedAt,
			&i.Product.UpdatedAt,
			&i.ItemTotalCost,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCartValueOfAllUsers = `-- name: GetCartValueOfAllUsers :many
SELECT c.user_id, SUM(ci.quantity * ci.price_at_time) AS total_cart_value
FROM carts c
JOIN cart_items ci ON c.id = ci.cart_id
GROUP BY c.user_id
ORDER BY total_cart_value
`

type GetCartValueOfAllUsersRow struct {
	UserID         uuid.UUID
	TotalCartValue int64
}

func (q *Queries) GetCartValueOfAllUsers(ctx context.Context) ([]GetCartValueOfAllUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getCartValueOfAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCartValueOfAllUsersRow
	for rows.Next() {
		var i GetCartValueOfAllUsersRow
		if err := rows.Scan(&i.UserID, &i.TotalCartValue); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getItemInCart = `-- name: GetItemInCart :one
SELECT ci.cart_id, ci.product_id, ci.quantity, ci.price_at_time, ci.created_at, ci.updated_at, p.id, p.title, p.supplier, p.category, p.price, p.image_url, p.description, p.product_location, p.created_at, p.updated_at 
FROM cart_items ci
JOIN products p ON ci.product_id = p.id
WHERE ci.cart_id = $1 AND ci.product_id = $2
`

type GetItemInCartParams struct {
	CartID    uuid.UUID
	ProductID uuid.UUID
}

type GetItemInCartRow struct {
	CartItem CartItem
	Product  Product
}

func (q *Queries) GetItemInCart(ctx context.Context, arg GetItemInCartParams) (GetItemInCartRow, error) {
	row := q.db.QueryRowContext(ctx, getItemInCart, arg.CartID, arg.ProductID)
	var i GetItemInCartRow
	err := row.Scan(
		&i.CartItem.CartID,
		&i.CartItem.ProductID,
		&i.CartItem.Quantity,
		&i.CartItem.PriceAtTime,
		&i.CartItem.CreatedAt,
		&i.CartItem.UpdatedAt,
		&i.Product.ID,
		&i.Product.Title,
		&i.Product.Supplier,
		&i.Product.Category,
		&i.Product.Price,
		&i.Product.ImageUrl,
		&i.Product.Description,
		&i.Product.ProductLocation,
		&i.Product.CreatedAt,
		&i.Product.UpdatedAt,
	)
	return i, err
}

const updateCartItem = `-- name: UpdateCartItem :exec
UPDATE cart_items
    SET quantity = $3,
        updated_at = NOW(),
        price_at_time = $4 
WHERE cart_id = $1 AND product_id = $2
`

type UpdateCartItemParams struct {
	CartID      uuid.UUID
	ProductID   uuid.UUID
	Quantity    int32
	PriceAtTime string
}

func (q *Queries) UpdateCartItem(ctx context.Context, arg UpdateCartItemParams) error {
	_, err := q.db.ExecContext(ctx, updateCartItem,
		arg.CartID,
		arg.ProductID,
		arg.Quantity,
		arg.PriceAtTime,
	)
	return err
}
