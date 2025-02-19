package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	app "github.com/alvarezjulia/meisterwerk-catalog/internal/application"
	productcmd "github.com/alvarezjulia/meisterwerk-catalog/internal/application/command"
	"github.com/alvarezjulia/meisterwerk-catalog/internal/application/product"
)

var (
	_ ServerInterface = &HTTPServer{}
)

type HTTPServer struct {
	app *app.App
}

func NewHTTPServer(a *app.App) *HTTPServer {
	return &HTTPServer{
		app: a,
	}
}

func (h *HTTPServer) CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var dtoProduct product.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&dtoProduct); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	domainProduct := dtoProduct.ToDomain()
	if domainProduct.Name == "" || domainProduct.Description == "" || domainProduct.Price == 0 {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := h.app.Commands.CreateProduct.Handle(ctx, &productcmd.CreateProductCommand{
		Name:        dtoProduct.Name,
		Description: dtoProduct.Description,
		Price:       dtoProduct.Price,
	})

	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *HTTPServer) GetProduct(w http.ResponseWriter, r *http.Request, id string) {
}

func (h *HTTPServer) ListProducts(w http.ResponseWriter, r *http.Request, params ListProductsParams) {
}

func (h *HTTPServer) UpdateProduct(w http.ResponseWriter, r *http.Request, id string) {
}

func (h *HTTPServer) DeleteProduct(w http.ResponseWriter, r *http.Request, ID string) {
	ctx := r.Context()
	if ID == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := h.app.Commands.DeleteProduct.Handle(ctx, &productcmd.DeleteProductCommand{
		ID,
	})
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPServer) BulkCreateProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var products []product.ProductResponse
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	productsList, err := mapBulkCreateProducts(products)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.app.Commands.BulkCreateProducts.Handle(ctx, &productcmd.BulkCreateProductsCommand{
		Products: productsList,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create products: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func mapBulkCreateProducts(products []product.ProductResponse) ([]productcmd.CreateProductCommand, error) {
	var productsList []productcmd.CreateProductCommand
	for _, product := range products {
		if product.Name == "" || product.Description == "" || product.Price == 0 {
			return nil, errors.New("invalid request payload")
		}
		productsList = append(productsList, productcmd.CreateProductCommand{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return productsList, nil
}

func (h *HTTPServer) BulkUpdateProducts(w http.ResponseWriter, r *http.Request) {
}

func (h *HTTPServer) BulkDeleteProducts(w http.ResponseWriter, r *http.Request) {
}
