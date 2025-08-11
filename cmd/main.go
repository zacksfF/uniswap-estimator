package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"uniswap-est/intrenal/config"
	"uniswap-est/intrenal/handlers"
	"uniswap-est/intrenal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize services
	blockchainService, err := services.NewBlockchainService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize blockchain service: %v", err)
	}
	defer blockchainService.Close()

	uniswapService := services.NewUniswapService(blockchainService)

	// Initialize handlers
	estimateHandler := handlers.NewEstimateHandler(uniswapService, cfg.RequestTimeout)
	healthHandler := handlers.NewHealthHandler(version)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Uniswap Estimator API",
		ServerHeader: "Uniswap-Estimator",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Routes
	setupRoutes(app, estimateHandler, healthHandler)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		app.Shutdown()
	}()

	// Start server
	serverAddr := cfg.GetServerAddress()
	log.Printf("Server starting on http://%s", serverAddr)
	log.Printf("Health check: http://%s/health", serverAddr)
	log.Printf("Estimate endpoint: http://%s/estimate", serverAddr)

	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes configures all application routes
func setupRoutes(app *fiber.App, estimateHandler *handlers.EstimateHandler, healthHandler *handlers.HealthHandler) {
	// API v1 routes
	v1 := app.Group("/api/v1")

	// Main endpoint - THE 1INCH REQUIREMENT!
	app.Get("/estimate", estimateHandler.EstimateSwap)
	v1.Get("/estimate", estimateHandler.EstimateSwap)

	// Health endpoints
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)
	v1.Get("/health", healthHandler.Health)
	v1.Get("/ready", healthHandler.Ready)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Uniswap V2 Estimator API",
			"version": version,
			"endpoints": map[string]string{
				"estimate": "/estimate?pool={pool}&src={src}&dst={dst}&src_amount={amount}",
				"health":   "/health",
				"ready":    "/ready",
			},
		})
	})
}

// customErrorHandler handles application-level errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"code":    code,
		"message": message,
	})
}
