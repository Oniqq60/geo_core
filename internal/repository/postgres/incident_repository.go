package postgres

import (
	"context"
	"errors"
	"fmt"
	"geo_system_core/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepository struct {
	db *pgxpool.Pool
}

func NewIncidentRepository(db *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) Create(ctx context.Context, req models.CreateIncidentRequest) (*models.Incident, error) {
	incident := &models.Incident{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		Severity:    req.Severity,
		Status:      req.Status,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if incident.Status == "" {
		incident.Status = "active"
	}

	query := `
		INSERT INTO incidents (id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		incident.ID, incident.Title, incident.Description,
		incident.Latitude, incident.Longitude, incident.Radius,
		incident.Severity, incident.Status, incident.IsActive,
		incident.CreatedAt, incident.UpdatedAt,
	).Scan(
		&incident.ID, &incident.Title, &incident.Description,
		&incident.Latitude, &incident.Longitude, &incident.Radius,
		&incident.Severity, &incident.Status, &incident.IsActive,
		&incident.CreatedAt, &incident.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create incident: %w", err)
	}

	return incident, nil
}

func (r *IncidentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Incident, error) {
	var incident models.Incident

	query := `
		SELECT id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at
		FROM incidents
		WHERE id = $1 AND is_active = true
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&incident.ID, &incident.Title, &incident.Description,
		&incident.Latitude, &incident.Longitude, &incident.Radius,
		&incident.Severity, &incident.Status, &incident.IsActive,
		&incident.CreatedAt, &incident.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("incident not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}

	return &incident, nil
}

func (r *IncidentRepository) List(ctx context.Context, page, limit int) ([]models.Incident, int, error) {
	offset := (page - 1) * limit

	// Получаем общее количество
	var total int
	countQuery := `SELECT COUNT(*) FROM incidents WHERE is_active = true`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count incidents: %w", err)
	}

	// Получаем список
	query := `
		SELECT id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at
		FROM incidents
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list incidents: %w", err)
	}
	defer rows.Close()

	var incidents []models.Incident
	for rows.Next() {
		var incident models.Incident
		err := rows.Scan(
			&incident.ID, &incident.Title, &incident.Description,
			&incident.Latitude, &incident.Longitude, &incident.Radius,
			&incident.Severity, &incident.Status, &incident.IsActive,
			&incident.CreatedAt, &incident.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan incident: %w", err)
		}
		incidents = append(incidents, incident)
	}

	return incidents, total, nil
}

func (r *IncidentRepository) Update(ctx context.Context, id uuid.UUID, req models.UpdateIncidentRequest) (*models.Incident, error) {
	// Сначала получаем текущий инцидент
	incident, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновляем поля
	if req.Title != nil {
		incident.Title = *req.Title
	}
	if req.Description != nil {
		incident.Description = *req.Description
	}
	if req.Latitude != nil {
		incident.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		incident.Longitude = *req.Longitude
	}
	if req.Radius != nil {
		incident.Radius = *req.Radius
	}
	if req.Severity != nil {
		incident.Severity = *req.Severity
	}
	if req.Status != nil {
		incident.Status = *req.Status
	}
	incident.UpdatedAt = time.Now()

	query := `
		UPDATE incidents
		SET title = $1, description = $2, latitude = $3, longitude = $4, radius = $5, severity = $6, status = $7, updated_at = $8
		WHERE id = $9 AND is_active = true
		RETURNING id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query,
		incident.Title, incident.Description,
		incident.Latitude, incident.Longitude, incident.Radius,
		incident.Severity, incident.Status, incident.UpdatedAt,
		id,
	).Scan(
		&incident.ID, &incident.Title, &incident.Description,
		&incident.Latitude, &incident.Longitude, &incident.Radius,
		&incident.Severity, &incident.Status, &incident.IsActive,
		&incident.CreatedAt, &incident.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update incident: %w", err)
	}

	return incident, nil
}

func (r *IncidentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE incidents SET is_active = false, updated_at = $1 WHERE id = $2 AND is_active = true`

	result, err := r.db.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete incident: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("incident not found")
	}

	return nil
}

func (r *IncidentRepository) FindNearby(ctx context.Context, lat, lng, maxDistance float64) ([]models.Incident, error) {
	// Используем формулу гаверсинуса для расчета расстояния
	// Используем подзапрос для фильтрации по расстоянию
	query := `
		SELECT id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at
		FROM (
			SELECT id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at,
			       6371000 * acos(
			           cos(radians($1)) * cos(radians(latitude)) *
			           cos(radians(longitude) - radians($2)) +
			           sin(radians($1)) * sin(radians(latitude))
			       ) AS distance
			FROM incidents
			WHERE is_active = true AND status = 'active'
		) AS incidents_with_distance
		WHERE distance <= $3
		ORDER BY distance
	`

	rows, err := r.db.Query(ctx, query, lat, lng, maxDistance)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby incidents: %w", err)
	}
	defer rows.Close()

	var incidents []models.Incident
	for rows.Next() {
		var incident models.Incident
		err := rows.Scan(
			&incident.ID, &incident.Title, &incident.Description,
			&incident.Latitude, &incident.Longitude, &incident.Radius,
			&incident.Severity, &incident.Status, &incident.IsActive,
			&incident.CreatedAt, &incident.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}

func (r *IncidentRepository) GetActiveIncidents(ctx context.Context) ([]models.Incident, error) {
	query := `
		SELECT id, title, description, latitude, longitude, radius, severity, status, is_active, created_at, updated_at
		FROM incidents
		WHERE is_active = true AND status = 'active'
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active incidents: %w", err)
	}
	defer rows.Close()

	var incidents []models.Incident
	for rows.Next() {
		var incident models.Incident
		err := rows.Scan(
			&incident.ID, &incident.Title, &incident.Description,
			&incident.Latitude, &incident.Longitude, &incident.Radius,
			&incident.Severity, &incident.Status, &incident.IsActive,
			&incident.CreatedAt, &incident.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		incidents = append(incidents, incident)
	}

	return incidents, nil
}
