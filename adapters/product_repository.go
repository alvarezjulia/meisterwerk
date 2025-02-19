package adapters

import (
	"context"
	"database/sql"
	"fmt"

	product "github.com/alvarezjulia/meisterwerk-catalog/internal/application/command"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *product.CreateProductCommand) error {
	query := `
	INSERT INTO products (name, description, price) 
	VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	query := `
	DELETE FROM products 
	WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (r *ProductRepository) CreateManyProducts(ctx context.Context, products []product.CreateProductCommand) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO products (name, description, price) VALUES ($1, $2, $3)`
	for _, product := range products {
		_, err := tx.Exec(query, product.Name, product.Description, product.Price)
		if err != nil {
			return fmt.Errorf("failed to create product: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

/*
func (r *ProductRepository) GetProducts(limit, offset int) ([]Product, error) {
	query := `SELECT id, name, description, price, created_at, updated_at FROM products LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) UpdateProduct(product *Product) error {
	query := `
	UPDATE products
	SET name = $1, description = $2, price = $3
	WHERE id = $4`
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *ProductRepository) UpdateManyProducts(products []Product) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4`
	for _, product := range products {
		_, err := tx.Exec(query, product.Name, product.Description, product.Price, product.ID)
		if err != nil {
			return fmt.Errorf("failed to update product %s: %w", product.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *ProductRepository) DeleteManyProducts(products []Product) error {
	query := `
	DELETE FROM products
	WHERE id = ANY($1)`
	_, err := r.db.Exec(query, products)
	if err != nil {
		return fmt.Errorf("failed to delete products: %w", err)
	}
	return nil
}
*/
