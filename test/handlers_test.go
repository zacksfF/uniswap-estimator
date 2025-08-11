package test

import (
	"net/http/httptest"
	"testing"
	"uniswap-est/intrenal/handlers"

	"github.com/gofiber/fiber/v2"
)

func TestHealthHandler(t *testing.T) {
	app := fiber.New()

	healthHandler := handlers.NewHealthHandler("test")
	app.Get("/health", healthHandler.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}
