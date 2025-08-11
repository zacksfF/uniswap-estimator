package handlers

import (
	"context"
	"net/http"
	"time"
	"uniswap-est/intrenal/models"
	"uniswap-est/intrenal/services"
	"uniswap-est/intrenal/utils"

	"github.com/gofiber/fiber/v2"
)

type EstimateHandler struct {
	uniswapService *services.UniswapService
	requestTimeout time.Duration
}

// NewEstimateHandler creates a new estimate handler
func NewEstimateHandler(uniswapService *services.UniswapService, timeout time.Duration) *EstimateHandler {
	return &EstimateHandler{
		uniswapService: uniswapService,
		requestTimeout: timeout,
	}
}

// EstimateSwap handles POST /estimate endpoint
// Example: GET /estimate?pool=0x...&src=0x...&dst=0x...&src_amount=1000000
func (h *EstimateHandler) EstimateSwap(c *fiber.Ctx) error {
	// Create request context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTimeout)
	defer cancel()

	// Parse query parameters into request model
	req := &models.EstimateRequest{
		Pool:      c.Query("pool"),
		Src:       c.Query("src"),
		Dst:       c.Query("dst"),
		SrcAmount: c.Query("src_amount"),
	}

	// Validate request
	if err := h.validateRequest(req); err != nil {
		return h.handleError(c, err)
	}

	// Process the estimate
	response, err := h.uniswapService.EstimateSwap(ctx, req)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return successful response
	return c.Status(http.StatusOK).JSON(response)
}

// validateRequest validates the incoming request
func (h *EstimateHandler) validateRequest(req *models.EstimateRequest) error {
	// Validate pool address
	if !utils.IsValidEthereumAddress(req.Pool) {
		return models.ErrInvalidPoolAddress
	}

	// Validate source token address
	if !utils.IsValidEthereumAddress(req.Src) {
		return models.NewAPIError(
			http.StatusBadRequest,
			"Invalid source token address",
			"Source token address must be a valid Ethereum address",
		)
	}

	// Validate destination token address
	if !utils.IsValidEthereumAddress(req.Dst) {
		return models.NewAPIError(
			http.StatusBadRequest,
			"Invalid destination token address",
			"Destination token address must be a valid Ethereum address",
		)
	}

	// Validate amount
	if !utils.IsValidAmount(req.SrcAmount) {
		return models.ErrInvalidAmount
	}

	// Ensure src and dst are different
	if utils.NormalizeAddress(req.Src) == utils.NormalizeAddress(req.Dst) {
		return models.NewAPIError(
			http.StatusBadRequest,
			"Invalid token pair",
			"Source and destination tokens must be different",
		)
	}

	return nil
}

// handleError handles API errors consistently
func (h *EstimateHandler) handleError(c *fiber.Ctx, err error) error {
	// Check if it's our custom API error
	if apiErr, ok := err.(*models.APIError); ok {
		return c.Status(apiErr.Code).JSON(apiErr)
	}

	// Handle context timeout
	if err == context.DeadlineExceeded {
		return c.Status(http.StatusRequestTimeout).JSON(&models.APIError{
			Code:    http.StatusRequestTimeout,
			Message: "Request timeout",
			Details: "The request took too long to process",
		})
	}

	// Generic server error
	return c.Status(http.StatusInternalServerError).JSON(&models.APIError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		Details: "An unexpected error occurred",
	})
}
