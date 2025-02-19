package productcmd

import (
	"context"
	"errors"
)

// BulkCreateProducts interface for handling bulk product creation
type BulkCreateProducts interface {
	Handle(ctx context.Context, cmd *BulkCreateProductsCommand) error
}

// BulkCreateProductsCommand represents a batch of products to be created
type BulkCreateProductsCommand struct {
	Products []CreateProductCommand `json:"products"`
}

//go:generate mockgen -destination=mocks/mock_bulk_create_products_handler.go -package=mocks -source=bulk_create_products.go CreateManyProducts
type BulkCreateProductsRepository interface {
	CreateManyProducts(ctx context.Context, products []CreateProductCommand) error
}

// bulkCreateProductsHandler implements BulkCreateProducts
type bulkCreateProductsHandler struct {
	repository BulkCreateProductsRepository
}

// NewBulkCreateProductsHandler creates a new instance
func NewBulkCreateProductsHandler(repository BulkCreateProductsRepository) *bulkCreateProductsHandler {
	return &bulkCreateProductsHandler{repository: repository}
}

// Handle processes bulk product creation internally
func (h *bulkCreateProductsHandler) Handle(ctx context.Context, cmd *BulkCreateProductsCommand) error {
	if len(cmd.Products) == 0 {
		return errors.New("no products provided")
	}
	return h.repository.CreateManyProducts(ctx, cmd.Products)
}
