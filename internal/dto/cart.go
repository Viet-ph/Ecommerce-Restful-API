package dto

import (
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated-at"`
}

type CartItem struct {
	ID          uuid.UUID `json:"id"`
	CartID      uuid.UUID `json:"cart_id"`
	ProductID   uuid.UUID `json:"product_id"`
	Quantity    int32     `json:"quantity"`
	PriceAtTime string    `json:"price_at_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func DbCartToDto(dbCart *db.Cart) Cart {
	return Cart{
		ID:        dbCart.ID,
		UserID:    dbCart.UserID,
		CreatedAt: dbCart.CreatedAt,
		UpdatedAt: dbCart.UpdatedAt,
	}
}

func DbCartItemToDto(dbCartItem *db.CartItem) CartItem {
	return CartItem{
		ID:          dbCartItem.ID,
		CartID:      dbCartItem.CartID,
		ProductID:   dbCartItem.ProductID,
		Quantity:    dbCartItem.Quantity,
		PriceAtTime: dbCartItem.PriceAtTime,
		CreatedAt:   dbCartItem.CreatedAt,
		UpdatedAt:   dbCartItem.UpdatedAt,
	}
}
