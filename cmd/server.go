package cmd

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alvarezjulia/meisterwerk-catalog/adapters"
	"github.com/alvarezjulia/meisterwerk-catalog/config"
	app "github.com/alvarezjulia/meisterwerk-catalog/internal/application"
	product "github.com/alvarezjulia/meisterwerk-catalog/internal/application/command"
	"github.com/alvarezjulia/meisterwerk-catalog/internal/infrastructure/cache"
	"github.com/alvarezjulia/meisterwerk-catalog/internal/infrastructure/db"
	"github.com/alvarezjulia/meisterwerk-catalog/migrations"
	"github.com/alvarezjulia/meisterwerk-catalog/ports"
)

type Dependencies struct {
	DB    *sql.DB
	Cache cache.Repository
}

func (d *Dependencies) Close() {
	if err := d.DB.Close(); err != nil {
		slog.Error("Error closing DB connection: %v", "error", err)
	}
	if err := d.Cache.Close(); err != nil {
		slog.Error("Error closing Cache connection: %v", "error", err)
	}
}

func StartDependencies(config *config.Config) *Dependencies {
	database := db.ConnectDB(config.DatabaseURL)
	cache, err := cache.New(config.RedisURL)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}

	return &Dependencies{DB: database, Cache: cache}
}

func Execute() {
	config, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}
	deps := StartDependencies(config)
	defer deps.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	migrationsPath := filepath.Join("migrations")
	if err := migrations.RunMigrations(deps.DB, migrationsPath); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		panic(err)
	}

	app := newApp(deps)
	RunServer(app, deps, config)
}

func newApp(deps *Dependencies) *app.App {
	productRepository := adapters.NewProductRepository(deps.DB)

	return &app.App{
		Commands: &app.Commands{
			CreateProduct: product.NewCreateProductHandler(productRepository),
			//UpdateProduct:      product.NewUpdateProductHandler(productRepository),
			DeleteProduct:      product.NewDeleteProductHandler(productRepository),
			BulkCreateProducts: product.NewBulkCreateProductsHandler(productRepository),
			//BulkUpdateProducts: product.NewBulkUpdateHandler(productRepository),
			//BulkDeleteProducts: product.NewBulkDeleteHandler(productRepository),
		},
		Queries: &app.Queries{
			//GetProducts:    query.NewGetProductsHandler(productRepository),
			//GetProductByID: query.NewGetProductByIDHandler(productRepository),
		},
	}
}

func RunServer(app *app.App, deps *Dependencies, config *config.Config) {
	mux := http.NewServeMux()
	httpServer := ports.NewHTTPServer(app)

	//mux.HandleFunc("/products", httpServer.GetProducts)
	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			httpServer.CreateProduct(w, r)
		}
	})
	mux.HandleFunc("/product/{id}", func(w http.ResponseWriter, r *http.Request) {
		ID := r.PathValue("id")

		switch r.Method {
		case http.MethodGet:
			httpServer.GetProduct(w, r, ID)
		case http.MethodPut:
			httpServer.UpdateProduct(w, r, ID)
		case http.MethodDelete:
			httpServer.DeleteProduct(w, r, ID)
		}
	})
	mux.HandleFunc("/products/bulk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			httpServer.BulkCreateProducts(w, r)
		case http.MethodPut:
			httpServer.BulkUpdateProducts(w, r)
		case http.MethodDelete:
			httpServer.BulkDeleteProducts(w, r)
		}
	})

	// protectedMux := middleware.Middleware(mux)

	slog.Info(fmt.Sprintf("Starting server on %s", config.Port))
	http.ListenAndServe(":"+config.Port, mux)
}
