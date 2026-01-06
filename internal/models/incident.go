package models

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Latitude    float64   `json:"latitude" db:"latitude"`
	Longitude   float64   `json:"longitude" db:"longitude"`
	Radius      float64   `json:"radius" db:"radius"`     // радиус в метрах
	Severity    string    `json:"severity" db:"severity"` // low, medium, high, critical
	Status      string    `json:"status" db:"status"`     // active, resolved
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateIncidentRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Radius      float64 `json:"radius" binding:"required,gt=0"`
	Severity    string  `json:"severity" binding:"required,oneof=low medium high critical"`
	Status      string  `json:"status" binding:"oneof=active resolved"`
}

type UpdateIncidentRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Radius      *float64 `json:"radius"`
	Severity    *string  `json:"severity" binding:"omitempty,oneof=low medium high critical"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active resolved"`
}

type IncidentResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Radius      float64   `json:"radius"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaginationParams struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}
