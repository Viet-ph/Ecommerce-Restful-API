package handler

import (
	"log"
	"net/http"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
	"github.com/google/uuid"
)

type CartHandler struct {
	CartService *service.CartService
}

func NewCartHandler(c *service.CartService) *CartHandler {
	return &CartHandler{
		CartService: c,
	}
}

func (c *CartHandler) AddNewCartItem() http.HandlerFunc {
	type request struct {
		ProductId    string `json:"product_id"`
		Quantity     int    `json:"quantity"`
		ProductPrice string `json:"product_price"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}
		cart, ok := r.Context().Value(middleware.ContextCartKey).(db.Cart)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not Cart type")
			return
		}

		productId, err := uuid.Parse(req.ProductId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Cannot parse product Id: "+err.Error())
			return
		}

		item, err := c.CartService.AddNewItemToCart(r.Context(), req.ProductPrice, int32(req.Quantity), cart.ID, productId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Cannot add item to cart: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, item)
	}
}

func (c *CartHandler) GetAllCartItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cart, ok := r.Context().Value(middleware.ContextCartKey).(db.Cart)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not Cart type")
			return
		}

		items, err := c.CartService.GetAllItemsInCart(r.Context(), cart.ID)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Cannot get items in cart: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, items)
	}
}

func (c *CartHandler) UpdateCartItem() http.HandlerFunc {
	type request struct {
		ProductId string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		cart, ok := r.Context().Value(middleware.ContextCartKey).(db.Cart)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not Cart type")
			return
		}

		productId, err := uuid.Parse(req.ProductId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "error parsing product Id: "+err.Error())
			return
		}

		err = c.CartService.RemoveOrUpdateItem(r.Context(), req.Quantity, cart.ID, productId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "error updating item: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, "Item updated")
	}
}

func (c *CartHandler) RemoveCartItem() http.HandlerFunc {
	type request struct {
		ProductId string `json:"product_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		cart, ok := r.Context().Value(middleware.ContextCartKey).(db.Cart)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not Cart type")
			return
		}

		productId, err := uuid.Parse(req.ProductId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "error parsing product Id: "+err.Error())
			return
		}

		err = c.CartService.RemoveCartItem(r.Context(), cart.ID, productId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "error updating item: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, "Item removed")
	}
}

func (c *CartHandler) UpdateCart() http.HandlerFunc {
	type request struct {
		ProductQuantityKV map[string]int `json:"product_quantity_kv"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		cart, ok := r.Context().Value(middleware.ContextCartKey).(db.Cart)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not Cart type")
			return
		}

		err = c.CartService.UpdateCart(r.Context(), cart.ID, req.ProductQuantityKV)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "error updating cart: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, "Cart updated")
	}
}

func (c *CartHandler) ClearCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cart, ok := r.Context().Value(middleware.ContextCartKey).(db.Cart)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not Cart type")
			return
		}

		err := c.CartService.ClearCart(r.Context(), cart.ID)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "error clearing cart: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, "Cart cleared")
	}
}
