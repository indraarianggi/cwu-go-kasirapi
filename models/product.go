package models

import "time"

type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`

	CategoryID *int      `json:"category_id,omitempty"`
	Category   *Category `json:"category,omitempty"`
}
