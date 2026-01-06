package handler

import (
	"geo_system_core/internal/models"
	"geo_system_core/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IncidentHandler struct {
	service *service.IncidentService
}

func NewIncidentHandler(service *service.IncidentService) *IncidentHandler {
	return &IncidentHandler{service: service}
}

func (h *IncidentHandler) Create(c *gin.Context) {
	var req models.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incident, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toIncidentResponse(incident))
}

func (h *IncidentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	incident, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "incident not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toIncidentResponse(incident))
}

func (h *IncidentHandler) List(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	incidents, total, err := h.service.List(c.Request.Context(), params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + params.Limit - 1) / params.Limit
	responses := make([]models.IncidentResponse, len(incidents))
	for i, incident := range incidents {
		responses[i] = toIncidentResponse(&incident)
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:       responses,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *IncidentHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incident, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "incident not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toIncidentResponse(incident))
}

func (h *IncidentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "incident not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "incident deactivated successfully"})
}

func toIncidentResponse(incident *models.Incident) models.IncidentResponse {
	return models.IncidentResponse{
		ID:          incident.ID,
		Title:       incident.Title,
		Description: incident.Description,
		Latitude:    incident.Latitude,
		Longitude:   incident.Longitude,
		Radius:      incident.Radius,
		Severity:    incident.Severity,
		Status:      incident.Status,
		IsActive:    incident.IsActive,
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
	}
}
