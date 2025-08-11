package utils

import (
	"errors"
	"math/big"
)

var (
	// Uniswap V2 constants for performance
	big997  = big.NewInt(997)  // Fee factor (1000 - 3)
	big1000 = big.NewInt(1000) // Fee denominator
	bigZero = big.NewInt(0)
)

// CalculateAmountOut implements Uniswap V2 math with 0.3% fee
// Formula: amountOut = (amountIn * 997 * reserveOut) / (reserveIn * 1000 + amountIn * 997)
// This is the CORE function that 1inch wants optimized!
func CalculateAmountOut(amountIn, reserveIn, reserveOut *big.Int) (*big.Int, error) {
	// Input validation
	if amountIn.Cmp(bigZero) <= 0 {
		return nil, errors.New("amount in must be positive")
	}
	if reserveIn.Cmp(bigZero) <= 0 || reserveOut.Cmp(bigZero) <= 0 {
		return nil, errors.New("insufficient liquidity")
	}

	// Performance-oriented calculation with minimal allocations
	// Using pre-allocated big integers to avoid memory allocations

	// amountInWithFee = amountIn * 997
	amountInWithFee := new(big.Int).Mul(amountIn, big997)

	// numerator = amountInWithFee * reserveOut
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)

	// denominator = reserveIn * 1000 + amountInWithFee
	denominator := new(big.Int).Mul(reserveIn, big1000)
	denominator.Add(denominator, amountInWithFee)

	// Final division
	amountOut := new(big.Int).Div(numerator, denominator)

	return amountOut, nil
}

// ConvertToTokenUnits converts amount considering token decimals
func ConvertToTokenUnits(amount *big.Int, decimals uint8) *big.Int {
	if decimals == 0 {
		return new(big.Int).Set(amount)
	}

	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return new(big.Int).Mul(amount, multiplier)
}

// ConvertFromTokenUnits converts amount from token units to base units
func ConvertFromTokenUnits(amount *big.Int, decimals uint8) *big.Int {
	if decimals == 0 {
		return new(big.Int).Set(amount)
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return new(big.Int).Div(amount, divisor)
}

// ParseBigInt safely parses string to big.Int
func ParseBigInt(s string) (*big.Int, error) {
	result := new(big.Int)
	_, ok := result.SetString(s, 10)
	if !ok {
		return nil, errors.New("invalid number format")
	}
	return result, nil
}
