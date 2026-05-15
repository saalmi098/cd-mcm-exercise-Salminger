package model

import "fmt"

// Product represents a product in the catalog.
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Validate checks whether the product has valid fields.
func (p *Product) Validate() bool {
	if p.Name == "" {
		return false
	}
	if p.Price < 0 {
		return false
	}
	return true
}

// Discount returns the price after applying a percentage discount.
func (p *Product) Discount(pct float64) float64 {
	if pct < 0 || pct > 100 {
		return p.Price
	}
	return p.Price * (1 - pct/100)
}

// IsExpensive reports whether the product price exceeds the given threshold.
func (p *Product) IsExpensive(threshold float64) bool {
	return p.Price > threshold
}

// String returns a human-readable representation of the product.
func (p *Product) String() string {
	return fmt.Sprintf("Product{ID: %d, Name: %q, Price: %.2f}", p.ID, p.Name, p.Price)
}
