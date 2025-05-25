package models

import (
	"time"
)

// Product represents the product model in the database
// @Description Product information
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`
	SellerID    uint      `json:"seller_id" binding:"required" example:"1"`
	Title       string    `json:"title" binding:"required" example:"iPhone 13 Pro"`
	Description string    `json:"description" example:"Latest iPhone model with pro camera system"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-01-01T00:00:00Z"`
} 