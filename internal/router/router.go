package router

import (
	"context"
	"geo_system_core/internal/config"
	"geo_system_core/internal/handler"
	"geo_system_core/internal/middleware"
	"geo_system_core/internal/repository/postgres"
	"geo_system_core/internal/repository/redis"
	"geo_system_core/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	incidentRepo *postgres.IncidentRepository,
	locationRepo *postgres.LocationRepository,
	queueRepo *redis.QueueRepository,
) *gin.Engine {
	// Инициализация сервисов
	incidentService := service.NewIncidentService(incidentRepo)
	locationService := service.NewLocationService(incidentRepo, locationRepo, queueRepo)
	statsService := service.NewStatsService(locationRepo, cfg.Stats.TimeWindowMinutes)
	webhookService := service.NewWebhookService(queueRepo, &cfg.Webhook)

	// Запускаем worker для обработки вебхуков
	ctx := context.Background()
	go webhookService.StartWorker(ctx)

	// Инициализация handlers
	incidentHandler := handler.NewIncidentHandler(incidentService)
	locationHandler := handler.NewLocationHandler(locationService)
	statsHandler := handler.NewStatsHandler(statsService)
	healthHandler := handler.NewHealthHandler()

	// Настройка роутера
	r := gin.Default()

	// Health check (публичный)
	r.GET("/api/v1/system/health", healthHandler.Health)

	// Публичный эндпоинт для проверки координат
	r.POST("/api/v1/location/check", locationHandler.Check)

	// Статистика (публичный)
	r.GET("/api/v1/incidents/stats", statsHandler.GetStats)

	// API для управления инцидентами (требует API-key)
	api := r.Group("/api/v1/incidents")
	api.Use(middleware.APIKeyAuth(cfg.Auth.APIKey))
	{
		api.POST("", incidentHandler.Create)
		api.GET("", incidentHandler.List)
		api.GET("/:id", incidentHandler.GetByID)
		api.PUT("/:id", incidentHandler.Update)
		api.DELETE("/:id", incidentHandler.Delete)
	}

	return r
}
