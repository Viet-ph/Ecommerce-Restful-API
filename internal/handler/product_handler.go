package handler

import (
	"log"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
	"github.com/google/uuid"
)

type ProductHandler struct {
	*service.ProductService
}

func NewProductHandler(p *service.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: p,
	}
}

func (p *ProductHandler) AddNewProduct() http.HandlerFunc {
	type request struct {
		Title           string `json:"title"`
		Supplier        string `json:"supplier"`
		Category        string `json:"category"`
		Price           string `json:"price"`
		ImageUrl        string `json:"image_url"`
		Description     string `json:"description"`
		ProductLocation string `json:"product_location"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		newProduct, err := p.Create(r.Context(),
			req.Title,
			req.Supplier,
			req.Category,
			req.Price,
			req.ImageUrl,
			req.Description,
			req.ProductLocation,
		)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error creating new product: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusCreated, newProduct)

	}
}

func (p *ProductHandler) GetProductDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId := r.PathValue("id")
		if productId == "" {
			helper.RespondWithError(w, http.StatusBadRequest, "Empty product Id")
			return
		}

		parsedId, err := uuid.Parse(productId)
		if err != nil {
			helper.RespondWithError(w, http.StatusBadRequest, "Error parsing UUID: "+err.Error())
			return
		}

		product, err := p.GetProductById(r.Context(), parsedId)
		if err != nil {
			helper.RespondWithError(w, http.StatusNotFound, "Couldn't find product with given id: "+productId)
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, product)
	}
}

func (p *ProductHandler) GetProductsWithFilters() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		products, err := p.ListProductsWithFilter(r.Context(), queryParams)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error getting products: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, products)
	}
}

func (p *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId := r.PathValue("id")
		if productId == "" {
			helper.RespondWithError(w, http.StatusBadRequest, "Empty product Id")
			return
		}

		parsedId, err := uuid.Parse(productId)
		if err != nil {
			helper.RespondWithError(w, http.StatusBadRequest, "Error parsing UUID: "+err.Error())
			return
		}

		err = p.DeleteProductById(r.Context(), parsedId)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete product with given id: "+productId)
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, "Product with id: "+productId+" deleted successfully.")
	}
}
