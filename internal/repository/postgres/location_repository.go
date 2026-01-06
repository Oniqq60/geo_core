package postgres

import (
	"context"
	"fmt"
	"geo_system_core/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) SaveCheck(ctx context.Context, userID string, lat, lng float64, hasDanger bool) error {
	query := `
		INSERT INTO location_checks (id, user_id, latitude, longitude, has_danger, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, uuid.New(), userID, lat, lng, hasDanger, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save location check: %w", err)
	}

	return nil
}

func (r *LocationRepository) GetZoneStats(ctx context.Context, timeWindowMinutes int) ([]models.ZoneStats, error) {
	// Используем параметризованный запрос для безопасности
	query := `
		SELECT 
			i.id as incident_id,
			i.title,
			COUNT(DISTINCT lc.user_id) as user_count
		FROM incidents i
		INNER JOIN location_checks lc ON 
			6371000 * acos(
				cos(radians(i.latitude)) * cos(radians(lc.latitude)) *
				cos(radians(lc.longitude) - radians(i.longitude)) +
				sin(radians(i.latitude)) * sin(radians(lc.latitude))
			) <= i.radius
		WHERE 
			i.is_active = true 
			AND i.status = 'active'
			AND lc.has_danger = true
			AND lc.created_at >= NOW() - (INTERVAL '1 minute' * $1)
		GROUP BY i.id, i.title
		ORDER BY user_count DESC
	`

	rows, err := r.db.Query(ctx, query, timeWindowMinutes)
	if err != nil {
		return nil, fmt.Errorf("failed to get zone stats: %w", err)
	}
	defer rows.Close()

	var stats []models.ZoneStats
	for rows.Next() {
		var stat models.ZoneStats
		err := rows.Scan(&stat.IncidentID, &stat.Title, &stat.UserCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan zone stat: %w", err)
		}
		stats = append(stats, stat)
	}

	return stats, nil
}
