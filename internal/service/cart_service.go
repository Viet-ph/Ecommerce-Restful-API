package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/dto"
	"github.com/google/uuid"
)

type CartService struct {
	queries *db.Queries
}

func NewCartService(q *db.Queries) *CartService {
	return &CartService{
		queries: q,
	}
}

func (cartService *CartService) CreateCart(ctx context.Context, userId uuid.UUID) (dto.Cart, error) {
	dbCart, err := cartService.queries.CreateCartForUser(ctx, db.CreateCartForUserParams{
		ID:        uuid.New(),
		UserID:    userId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return dto.Cart{}, fmt.Errorf("unable to create new cart: %v", err)
	}

	return dto.DbCartToDto(&dbCart), nil
}

func (cartService *CartService) AddNewItemToCart(ctx context.Context, productPrice string, quantity int32, cartId, productId uuid.UUID) (dto.CartItem, error) {
	dbCartItem, err := cartService.queries.AddNewItemToCart(ctx, db.AddNewItemToCartParams{
		CartID:      cartId,
		ProductID:   productId,
		Quantity:    quantity,
		PriceAtTime: productPrice,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return dto.CartItem{}, fmt.Errorf("unable to add new item into cart: %v", err)
	}

	return dto.DbCartItemToDto(&dbCartItem), nil
}

func (cartService *CartService) GetAllItemsInCart(ctx context.Context, cartId uuid.UUID) ([]dto.CartItemProductJoined, error) {
	dbCartItemsRow, err := cartService.queries.GetAllItemsInCart(ctx, cartId)
	if err != nil {
		return []dto.CartItemProductJoined{}, fmt.Errorf("unable to get all items in cart: %v", err)
	}

	itemProductJoined := make([]dto.CartItemProductJoined, 0, len(dbCartItemsRow))
	for _, itemRow := range dbCartItemsRow {
		itemProductJoined = append(itemProductJoined, dto.CartItemProductJoined{
			Item:          dto.DbCartItemToDto(&itemRow.CartItem),
			Product:       dto.DbProductToDto(&itemRow.Product),
			ItemTotalCost: itemRow.ItemTotalCost,
		})
	}

	return itemProductJoined, nil
}

func (cartService *CartService) RemoveOrUpdateItem(ctx context.Context, quantity int, cartId, productId uuid.UUID) error {
	if quantity == 0 {
		return cartService.RemoveCartItem(ctx, cartId, productId)
	}

	//Check if the item exist in cart
	itemProductJoined, err := cartService.queries.GetItemInCart(ctx, db.GetItemInCartParams{
		CartID:    cartId,
		ProductID: productId,
	})

	if err == sql.ErrNoRows {
		_, err = cartService.AddNewItemToCart(ctx,
			itemProductJoined.Product.Price,
			int32(quantity),
			cartId,
			productId)
		if err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("unable to check if item is already exist in cart: %v", err)
	} else {
		err = cartService.queries.UpdateCartItem(ctx, db.UpdateCartItemParams{
			CartID:      itemProductJoined.CartItem.CartID,
			ProductID:   itemProductJoined.Product.ID,
			Quantity:    int32(quantity),
			PriceAtTime: itemProductJoined.Product.Price,
		})
		if err != nil {
			return fmt.Errorf("unable to update cart item: %v", err)
		}
	}
	return nil
}

func (cartService *CartService) RemoveCartItem(ctx context.Context, cartId, productId uuid.UUID) error {
	err := cartService.queries.DeleteItemInCart(ctx, db.DeleteItemInCartParams{
		CartID:    cartId,
		ProductID: productId,
	})
	if err != nil {
		return fmt.Errorf("unable to remove item in cart: %v", err)
	}

	return nil
}

func (cartService *CartService) ClearCart(ctx context.Context, userId uuid.UUID) error {
	err := cartService.queries.DeleteAllItemsInCart(ctx, userId)
	if err != nil {
		return fmt.Errorf("unable clear cart: %v", err)
	}

	return nil
}

func (cartService *CartService) UpdateCart(ctx context.Context, cartId uuid.UUID, productIdAndQuantity map[uuid.UUID]int) error {
	for productId, quantity := range productIdAndQuantity {
		err := cartService.RemoveOrUpdateItem(ctx, quantity, cartId, productId)
		if err != nil {
			return fmt.Errorf("unable to update item: %v, error meaasge: %v", productId, err)
		}
	}

	return nil
}
