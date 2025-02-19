package app

import product "github.com/alvarezjulia/meisterwerk-catalog/internal/application/command"

type App struct {
	Commands *Commands
	Queries  *Queries
}

type Commands struct {
	CreateProduct product.CreateProduct
	// UpdateProduct      product.UpdateProduct
	DeleteProduct      product.DeleteProduct
	BulkCreateProducts product.BulkCreateProducts
	// BulkUpdateProducts product.BulkUpdateProducts
	// BulkDeleteProducts product.BulkDeleteProducts
}

type Queries struct {
	// GetProducts    product.GetProducts
	// GetProductByID product.GetProductByID
}
