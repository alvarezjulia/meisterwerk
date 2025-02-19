package productcmd

import (
	"context"
	"errors"
)

type DeleteProduct interface {
	Handle(ctx context.Context, cmd *DeleteProductCommand) error
}

type DeleteProductCommand struct {
	ID string
}

//go:generate mockgen -destination=mocks/mock_delete_product_handler.go -package=mocks -source=delete_product.go DeleteProduct
type DeleteProductRepository interface {
	DeleteProduct(ctx context.Context, id string) error
}

type deleteProductHandler struct {
	repository DeleteProductRepository
}

// NewDeleteProductHandler creates a new instance of DeleteProductHandler
func NewDeleteProductHandler(repository DeleteProductRepository) *deleteProductHandler {
	return &deleteProductHandler{repository: repository}
}

// Handle executes the delete product command
func (h *deleteProductHandler) Handle(ctx context.Context, cmd *DeleteProductCommand) error {
	if cmd.ID == "" {
		return errors.New("product ID is required")
	}

	return h.repository.DeleteProduct(ctx, cmd.ID)
}
