package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/dto"
	"github.com/google/uuid"
)

type ProductService struct {
	queries *db.Queries
}

func NewProductService(q *db.Queries) *ProductService {
	return &ProductService{
		queries: q,
	}
}

func (productService *ProductService) Create(
	ctx context.Context,
	title,
	supplier,
	category,
	price,
	imageUrl,
	description,
	productLocation string,
) (dto.Product, error) {
	product, err := productService.queries.CreateProduct(ctx, db.CreateProductParams{
		ID:              uuid.New(),
		Title:           title,
		Supplier:        supplier,
		Category:        category,
		Price:           price,
		ImageUrl:        imageUrl,
		Description:     description,
		ProductLocation: productLocation,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	})
	if err != nil {
		return dto.Product{}, fmt.Errorf("unable to create new product: %v", err)
	}

	return dto.DbProductToDto(&product), nil
}

func (productService *ProductService) GetProductById(ctx context.Context, id uuid.UUID) (dto.Product, error) {
	product, err := productService.queries.GetProductById(ctx, id)
	if err != nil {
		return dto.Product{}, fmt.Errorf("unable to get product: %v", err)
	}

	return dto.DbProductToDto(&product), nil
}

func (productService *ProductService) ListProductsWithFilter(ctx context.Context, URLQuery url.Values) ([]dto.Product, error) {
	filterParams := db.ListProductsWithFilterParams{}
	for k := range URLQuery {
		if k == "category" {
			filterParams.Category.String = URLQuery.Get(k)
			filterParams.Category.Valid = true
		}
		if k == "supplier" {
			filterParams.Supplier.String = URLQuery.Get(k)
			filterParams.Supplier.Valid = true
		}
		if k == "productLocation" {
			filterParams.ProductLocation.String = URLQuery.Get(k)
			filterParams.ProductLocation.Valid = true
		}
		if k == "limit" {
			limitInt64, err := strconv.ParseInt(URLQuery.Get(k), 10, 32)
			if err != nil {
				return []dto.Product{}, fmt.Errorf("error getting products limit: %v", err)
			}
			filterParams.Lim.Int32 = int32(limitInt64)
			filterParams.Lim.Valid = true
		}
		if k == "orderBy" {
			filterParams.OrderBy.String = URLQuery.Get(k)
			filterParams.OrderBy.Valid = true
		}
	}
	dbProducts, err := productService.queries.ListProductsWithFilter(ctx, filterParams)
	if err != nil {
		return []dto.Product{}, fmt.Errorf("unable to get products with filter: %v", err)
	}

	products := make([]dto.Product, 0, len(dbProducts))
	for _, dbProduct := range dbProducts {
		products = append(products, dto.DbProductToDto(&dbProduct))
	}

	return products, nil
}

func (productService *ProductService) DeleteProductById(ctx context.Context, id uuid.UUID) error {
	err := productService.queries.DeleteProductById(ctx, id)
	if err != nil {
		return fmt.Errorf("unable to delete product: %v", err)
	}

	return nil
}
