package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"geo_system_core/internal/config"
	"geo_system_core/internal/models"
	"geo_system_core/internal/repository/redis"
	"io"
	"net/http"
	"time"
)

type WebhookService struct {
	queueRepo *redis.QueueRepository
	config    *config.WebhookConfig
	client    *http.Client
}

func NewWebhookService(queueRepo *redis.QueueRepository, cfg *config.WebhookConfig) *WebhookService {
	return &WebhookService{
		queueRepo: queueRepo,
		config:    cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (s *WebhookService) StartWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			payload, err := s.queueRepo.DequeueWebhook(ctx)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			if payload == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			// Отправляем вебхук с retry
			s.sendWebhookWithRetry(ctx, payload)
		}
	}
}

func (s *WebhookService) sendWebhookWithRetry(ctx context.Context, payload *models.WebhookPayload) {
	for attempt := 1; attempt <= s.config.RetryAttempts; attempt++ {
		err := s.sendWebhook(ctx, payload)
		if err == nil {
			return // Успешно отправлено
		}

		if attempt < s.config.RetryAttempts {
			// Экспоненциальная задержка
			delay := time.Duration(attempt) * s.config.RetryDelay
			time.Sleep(delay)
		}
	}
}

func (s *WebhookService) sendWebhook(ctx context.Context, payload *models.WebhookPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.URL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
