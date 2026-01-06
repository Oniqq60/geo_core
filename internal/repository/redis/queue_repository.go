package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"geo_system_core/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type QueueRepository struct {
	client *redis.Client
}

func NewQueueRepository(client *redis.Client) *QueueRepository {
	return &QueueRepository{client: client}
}

const (
	webhookQueueKey = "webhook:queue"
	cacheKeyPrefix  = "incidents:active"
)

func (r *QueueRepository) EnqueueWebhook(ctx context.Context, payload models.WebhookPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	err = r.client.LPush(ctx, webhookQueueKey, data).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue webhook: %w", err)
	}

	return nil
}

func (r *QueueRepository) DequeueWebhook(ctx context.Context) (*models.WebhookPayload, error) {
	result, err := r.client.BRPop(ctx, 5*time.Second, webhookQueueKey).Result()
	if err == redis.Nil {
		return nil, nil // Очередь пуста
	}
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue webhook: %w", err)
	}

	var payload models.WebhookPayload
	err = json.Unmarshal([]byte(result[1]), &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook payload: %w", err)
	}

	return &payload, nil
}

func (r *QueueRepository) CacheActiveIncidents(ctx context.Context, incidents []models.Incident, ttl time.Duration) error {
	data, err := json.Marshal(incidents)
	if err != nil {
		return fmt.Errorf("failed to marshal incidents: %w", err)
	}

	err = r.client.Set(ctx, cacheKeyPrefix, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache incidents: %w", err)
	}

	return nil
}

func (r *QueueRepository) GetCachedActiveIncidents(ctx context.Context) ([]models.Incident, error) {
	data, err := r.client.Get(ctx, cacheKeyPrefix).Result()
	if err == redis.Nil {
		return nil, nil // Кэш пуст
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cached incidents: %w", err)
	}

	var incidents []models.Incident
	err = json.Unmarshal([]byte(data), &incidents)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal incidents: %w", err)
	}

	return incidents, nil
}
