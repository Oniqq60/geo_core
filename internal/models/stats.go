package models

import "github.com/google/uuid"

type ZoneStats struct {
	IncidentID uuid.UUID `json:"incident_id"`
	Title      string    `json:"title"`
	UserCount  int       `json:"user_count"`
}

type StatsResponse struct {
	Zones []ZoneStats `json:"zones"`
	Total int         `json:"total"`
}
