package service

import (
	"context"
	"fmt"
	"geo_system_core/internal/models"
	"geo_system_core/internal/repository/postgres"
	"geo_system_core/internal/repository/redis"
	"math"
	"time"
)

type LocationService struct {
	incidentRepo *postgres.IncidentRepository
	locationRepo *postgres.LocationRepository
	queueRepo    *redis.QueueRepository
}

func NewLocationService(
	incidentRepo *postgres.IncidentRepository,
	locationRepo *postgres.LocationRepository,
	queueRepo *redis.QueueRepository,
) *LocationService {
	return &LocationService{
		incidentRepo: incidentRepo,
		locationRepo: locationRepo,
		queueRepo:    queueRepo,
	}
}

// CalculateDistance вычисляет расстояние между двумя точками в метрах (формула гаверсинуса)
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // радиус Земли в метрах

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (s *LocationService) CheckLocation(ctx context.Context, req models.LocationCheckRequest) (*models.LocationCheckResponse, error) {
	// Валидация координат
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, fmt.Errorf("invalid latitude: must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, fmt.Errorf("invalid longitude: must be between -180 and 180")
	}

	// Ищем ближайшие инциденты (в радиусе 10 км для оптимизации)
	maxSearchDistance := 10000.0 // 10 км
	incidents, err := s.incidentRepo.FindNearby(ctx, req.Latitude, req.Longitude, maxSearchDistance)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby incidents: %w", err)
	}

	var nearbyIncidents []models.NearbyIncident
	hasDanger := false

	for _, incident := range incidents {
		distance := CalculateDistance(req.Latitude, req.Longitude, incident.Latitude, incident.Longitude)

		// Проверяем, находится ли пользователь в радиусе опасности
		if distance <= incident.Radius {
			hasDanger = true
			nearbyIncidents = append(nearbyIncidents, models.NearbyIncident{
				ID:          incident.ID,
				Title:       incident.Title,
				Description: incident.Description,
				Latitude:    incident.Latitude,
				Longitude:   incident.Longitude,
				Radius:      incident.Radius,
				Severity:    incident.Severity,
				Distance:    distance,
			})
		}
	}

	response := &models.LocationCheckResponse{
		NearbyIncidents: nearbyIncidents,
		HasDanger:       hasDanger,
	}

	// Сохраняем факт проверки в БД (асинхронно через горутину)
	go func() {
		ctx := context.Background()
		_ = s.locationRepo.SaveCheck(ctx, req.UserID, req.Latitude, req.Longitude, hasDanger)
	}()

	// Если есть опасности, ставим задачу на отправку вебхука
	if hasDanger {
		go func() {
			ctx := context.Background()
			payload := models.WebhookPayload{
				UserID:    req.UserID,
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
				Timestamp: time.Now(),
				Incidents: nearbyIncidents,
			}
			_ = s.queueRepo.EnqueueWebhook(ctx, payload)
		}()
	}

	return response, nil
}
