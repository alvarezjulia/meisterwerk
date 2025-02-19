package product

import (
	"time"

	"github.com/alvarezjulia/meisterwerk-catalog/domain"
)

// ProductResponse represents the API response structure
type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	TaxRate     TaxDTO    `json:"tax_rate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaxDTO represents the API tax structure
type TaxDTO struct {
	ID          string  `json:"id"`
	Rate        float64 `json:"rate"`
	Description string  `json:"description"`
	Country     string  `json:"country"`
	Region      string  `json:"region"`
}

// ToDTO converts a domain product to a DTO
func ToDTO(p *domain.Product) *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		TaxRate: TaxDTO{
			ID:          p.TaxRate.ID,
			Rate:        p.TaxRate.Rate,
			Description: p.TaxRate.Description,
			Country:     p.TaxRate.Country,
			Region:      p.TaxRate.Region,
		},
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToDomain converts a DTO to a domain product
func (dto *ProductResponse) ToDomain() *domain.Product {
	return &domain.Product{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		TaxRate: domain.Tax{
			ID:          dto.TaxRate.ID,
			Rate:        dto.TaxRate.Rate,
			Description: dto.TaxRate.Description,
			Country:     dto.TaxRate.Country,
			Region:      dto.TaxRate.Region,
		},
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
