package models

import (
	"math/big"
)

// EstimateRequest represents the input parameters for swap estimation
type EstimateRequest struct {
	Pool      string `json:"pool" validate:"required,len=42"`      // Uniswap V2 pair address
	Src       string `json:"src" validate:"required,len=42"`       // Source token address
	Dst       string `json:"dst" validate:"required,len=42"`       // Destination token address
	SrcAmount string `json:"src_amount" validate:"required"`       // Input amount as string
}

// EstimateResponse represents the API response
type EstimateResponse struct {
	DstAmount string `json:"dst_amount"` // Output amount calculated off-chain
}

// TokenInfo holds token metadata
type TokenInfo struct {
	Address  string
	Decimals uint8
	Symbol   string
}

// PoolReserves holds current pool state from blockchain
type PoolReserves struct {
	Reserve0   *big.Int
	Reserve1   *big.Int
	Token0     string
	Token1     string
	BlockTime  uint32
}

// SwapCalculation holds intermediate calculation data
type SwapCalculation struct {
	AmountIn     *big.Int
	AmountOut    *big.Int
	ReserveIn    *big.Int
	ReserveOut   *big.Int
	TokenIn      *TokenInfo
	TokenOut     *TokenInfo
}