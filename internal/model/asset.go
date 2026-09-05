package model

import "time"

type AssetCategory struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Asset struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	CategoryID    int64     `json:"category_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	PurchasePrice float64   `json:"purchase_price"`
	CurrentValue  float64   `json:"current_value"`
	PurchaseDate  time.Time `json:"purchase_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Category *AssetCategory `json:"category,omitempty"`
}
