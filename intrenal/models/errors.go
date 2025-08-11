package models

import (
	"fmt"
	"net/http"
)

// Custom error types for better error handling
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

// Predefined errors
var (
	ErrInvalidPoolAddress = &APIError{
		Code:    http.StatusBadRequest,
		Message: "Invalid pool address",
		Details: "Pool address must be a valid Ethereum address (42 characters)",
	}
	
	ErrInvalidTokenAddress = &APIError{
		Code:    http.StatusBadRequest,
		Message: "Invalid token address",
		Details: "Token addresses must be valid Ethereum addresses",
	}
	
	ErrInvalidAmount = &APIError{
		Code:    http.StatusBadRequest,
		Message: "Invalid amount",
		Details: "Amount must be a positive integer",
	}
	
	ErrPoolNotFound = &APIError{
		Code:    http.StatusNotFound,
		Message: "Pool not found",
		Details: "The specified pool address does not exist or is not a Uniswap V2 pair",
	}
	
	ErrInsufficientLiquidity = &APIError{
		Code:    http.StatusBadRequest,
		Message: "Insufficient liquidity",
		Details: "The pool does not have enough liquidity for this swap",
	}
	
	ErrBlockchainConnection = &APIError{
		Code:    http.StatusServiceUnavailable,
		Message: "Blockchain connection error",
		Details: "Unable to fetch current state from Ethereum network",
	}
)

// NewAPIError creates a custom API error
func NewAPIError(code int, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}