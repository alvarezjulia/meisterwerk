package domain

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	TaxRate     Tax       `json:"tax_rate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Tax struct {
	ID          string    `json:"id"`
	Rate        float64   `json:"rate"`
	Description string    `json:"description"`
	Country     string    `json:"country"`
	Region      string    `json:"region"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
