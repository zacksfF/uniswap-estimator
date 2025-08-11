package test

import (
	"math/big"
	"testing"
	"uniswap-est/intrenal/utils"
)

// BenchmarkCalculateAmountOut - THE 1INCH REQUIREMENT!
func BenchmarkCalculateAmountOut(b *testing.B) {
	// Realistic Uniswap V2 pool values
	amountIn := big.NewInt(1000000000)       // 1000 USDT (6 decimals)
	reserveIn := big.NewInt(100000000000000) // 100M USDT reserve

	reserveOut := new(big.Int)
	reserveOut.SetString("50000000000000000000000", 10) // 50k ETH reserve

	b.ResetTimer()
	b.ReportAllocs() // Show allocations - 1inch cares about this!

	for i := 0; i < b.N; i++ {
		_, _ = utils.CalculateAmountOut(amountIn, reserveIn, reserveOut)
	}
}

// BenchmarkMultipleCalculations - Test performance under load
func BenchmarkMultipleCalculations(b *testing.B) {
	testCases := []*big.Int{
		big.NewInt(1000000),     // 1 USDT
		big.NewInt(100000000),   // 100 USDT
		big.NewInt(1000000000),  // 1000 USDT
		big.NewInt(10000000000), // 10k USDT
	}

	reserveIn := big.NewInt(100000000000000)
	reserveOut := new(big.Int)
	reserveOut.SetString("50000000000000000000000", 10)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, amount := range testCases {
			_, _ = utils.CalculateAmountOut(amount, reserveIn, reserveOut)
		}
	}
}
