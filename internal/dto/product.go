package dto

import (
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/google/uuid"
)

type Product struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Supplier        string    `json:"supplier"`
	Category        string    `json:"category"`
	Price           string    `json:"price"`
	ImageUrl        string    `json:"image_url"`
	Description     string    `json:"description"`
	ProductLocation string    `json:"product_location"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func DbProductToDto(dbProduct *db.Product) Product {
	return Product{
		ID:              dbProduct.ID,
		Title:           dbProduct.Title,
		Supplier:        dbProduct.Supplier,
		Category:        dbProduct.Category,
		Price:           dbProduct.Price,
		ImageUrl:        dbProduct.ImageUrl,
		Description:     dbProduct.Description,
		ProductLocation: dbProduct.ProductLocation,
		CreatedAt:       dbProduct.CreatedAt,
		UpdatedAt:       dbProduct.UpdatedAt,
	}
}
