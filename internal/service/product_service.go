package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/google/uuid"
)

type ProductService struct {
	queries     *db.Queries
}

func NewProductService(q *db.Queries) *ProductService {
	return &ProductService{
		queries: q,
	}
}

