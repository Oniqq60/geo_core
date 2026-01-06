package service

import (
	"context"
	"geo_system_core/internal/models"
	"geo_system_core/internal/repository/postgres"
)

type StatsService struct {
	locationRepo *postgres.LocationRepository
	timeWindow   int
}

func NewStatsService(locationRepo *postgres.LocationRepository, timeWindowMinutes int) *StatsService {
	return &StatsService{
		locationRepo: locationRepo,
		timeWindow:   timeWindowMinutes,
	}
}

func (s *StatsService) GetZoneStats(ctx context.Context) (*models.StatsResponse, error) {
	stats, err := s.locationRepo.GetZoneStats(ctx, s.timeWindow)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, stat := range stats {
		total += stat.UserCount
	}

	return &models.StatsResponse{
		Zones: stats,
		Total: total,
	}, nil
}
