package models

import (
	"time"

	"github.com/google/uuid"
)

type LocationCheckRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	UserID    string  `json:"user_id" binding:"required"`
}

type LocationCheckResponse struct {
	NearbyIncidents []NearbyIncident `json:"nearby_incidents"`
	HasDanger       bool             `json:"has_danger"`
}

type NearbyIncident struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Radius      float64   `json:"radius"`
	Severity    string    `json:"severity"`
	Distance    float64   `json:"distance"` // расстояние в метрах
}

type LocationCheckLog struct {
	ID        uuid.UUID `db:"id"`
	UserID    string    `db:"user_id"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	HasDanger bool      `db:"has_danger"`
	CreatedAt time.Time `db:"created_at"`
}

type WebhookPayload struct {
	UserID    string           `json:"user_id"`
	Latitude  float64          `json:"latitude"`
	Longitude float64          `json:"longitude"`
	Timestamp time.Time        `json:"timestamp"`
	Incidents []NearbyIncident `json:"incidents"`
}
