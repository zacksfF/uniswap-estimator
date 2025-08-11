package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
}

// Health handles GET /health endpoint
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	uptime := time.Since(h.startTime)
	
	response := HealthResponse{
		Status:    "healthy",
		Version:   h.version,
		Uptime:    uptime.String(),
		Timestamp: time.Now().UTC(),
	}

	return c.Status(http.StatusOK).JSON(response)
}

// Ready handles GET /ready endpoint (for Kubernetes readiness probe)
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// You can add more sophisticated readiness checks here
	// For example: check blockchain connection, database connection, etc.
	
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "ready",
		"timestamp": time.Now().UTC(),
	})
}