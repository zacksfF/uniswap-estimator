package test

import (
	"math/big"
	"testing"
	"uniswap-est/intrenal/utils"
)

func TestCalculateAmountOut(t *testing.T) {
	// Test case: 1000 USDT -> ETH
	amountIn := big.NewInt(1000000000)    // 1000 USDT (6 decimals)
	reserveIn := big.NewInt(100000000000) // 100k USDT reserve

	// Fix: Use SetString for large numbers
	reserveOut := new(big.Int)
	reserveOut.SetString("50000000000000000000", 10) // 50 ETH reserve (18 decimals)

	result, err := utils.CalculateAmountOut(amountIn, reserveIn, reserveOut)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Cmp(big.NewInt(0)) <= 0 {
		t.Fatalf("Expected positive result, got %v", result)
	}

	t.Logf("AmountOut: %s", result.String())
}

