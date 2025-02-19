package productcmd

import (
	"context"
)

type CreateProduct interface {
	Handle(ctx context.Context, cmd *CreateProductCommand) error
}

type CreateProductCommand struct {
	Name        string  `json:"product_name"`
	Description string  `json:"product_description"`
	Price       float64 `json:"product_price"`
}

//go:generate mockgen -destination=mocks/mock_create_product_handler.go -package=mocks -source=create_product.go CreateProduct
type CreateProductRepository interface {
	CreateProduct(ctx context.Context, product *CreateProductCommand) error
}

type createProductHandler struct {
	repository CreateProductRepository
}

func NewCreateProductHandler(repository CreateProductRepository) *createProductHandler {
	return &createProductHandler{repository: repository}
}

func (h *createProductHandler) Handle(ctx context.Context, cmd *CreateProductCommand) error {
	return h.repository.CreateProduct(ctx, cmd)
}
