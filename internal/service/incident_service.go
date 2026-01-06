package service

import (
	"context"
	"geo_system_core/internal/models"
	"geo_system_core/internal/repository/postgres"
)

type IncidentService struct {
	repo *postgres.IncidentRepository
}

func NewIncidentService(repo *postgres.IncidentRepository) *IncidentService {
	return &IncidentService{repo: repo}
}

func (s *IncidentService) Create(ctx context.Context, req models.CreateIncidentRequest) (*models.Incident, error) {
	return s.repo.Create(ctx, req)
}

func (s *IncidentService) GetByID(ctx context.Context, id string) (*models.Incident, error) {
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, uuid)
}

func (s *IncidentService) List(ctx context.Context, page, limit int) ([]models.Incident, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.List(ctx, page, limit)
}

func (s *IncidentService) Update(ctx context.Context, id string, req models.UpdateIncidentRequest) (*models.Incident, error) {
	uuid, err := parseUUID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, uuid, req)
}

func (s *IncidentService) Delete(ctx context.Context, id string) error {
	uuid, err := parseUUID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, uuid)
}

func (s *IncidentService) GetActiveIncidents(ctx context.Context) ([]models.Incident, error) {
	return s.repo.GetActiveIncidents(ctx)
}
