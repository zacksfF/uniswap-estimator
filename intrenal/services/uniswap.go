package services

import (
	"context"
	"math/big"
	"uniswap-est/intrenal/models"
	"uniswap-est/intrenal/utils"
)

type UniswapService struct {
	blockchain *BlockchainService
}

// NewUniswapService creates a new Uniswap service
func NewUniswapService(blockchain *BlockchainService) *UniswapService {
	return &UniswapService{
		blockchain: blockchain,
	}
}

// EstimateSwap performs the complete swap estimation
// This is the main function that orchestrates everything!
func (us *UniswapService) EstimateSwap(ctx context.Context, req *models.EstimateRequest) (*models.EstimateResponse, error) {
	// Step 1: Normalize addresses
	poolAddr := utils.NormalizeAddress(req.Pool)
	srcAddr := utils.NormalizeAddress(req.Src)
	dstAddr := utils.NormalizeAddress(req.Dst)

	// Step 2: Parse input amount
	amountIn, err := utils.ParseBigInt(req.SrcAmount)
	if err != nil {
		return nil, models.ErrInvalidAmount
	}

	// Step 3: Fetch pool reserves from blockchain (CURRENT STATE!)
	reserves, err := us.blockchain.GetPoolReserves(ctx, poolAddr)
	if err != nil {
		return nil, err
	}

	// Step 4: Get token information
	srcToken, err := us.blockchain.GetTokenInfo(ctx, srcAddr)
	if err != nil {
		return nil, err
	}

	dstToken, err := us.blockchain.GetTokenInfo(ctx, dstAddr)
	if err != nil {
		return nil, err
	}

	// Step 5: Determine token order and reserves
	calculation, err := us.setupCalculation(amountIn, reserves, srcToken, dstToken)
	if err != nil {
		return nil, err
	}

	// Step 6: Apply Uniswap V2 math (THE CORE CALCULATION!)
	amountOut, err := utils.CalculateAmountOut(
		calculation.AmountIn,
		calculation.ReserveIn,
		calculation.ReserveOut,
	)
	if err != nil {
		return nil, models.ErrInsufficientLiquidity
	}

	return &models.EstimateResponse{
		DstAmount: amountOut.String(),
	}, nil
}

// setupCalculation determines which reserves to use based on token order
func (us *UniswapService) setupCalculation(
	amountIn *big.Int,
	reserves *models.PoolReserves,
	srcToken, dstToken *models.TokenInfo,
) (*models.SwapCalculation, error) {

	srcAddr := utils.NormalizeAddress(srcToken.Address)
	dstAddr := utils.NormalizeAddress(dstToken.Address)

	// Determine if src is token0 or token1 in the pair
	var reserveIn, reserveOut *big.Int

	if srcAddr == reserves.Token0 && dstAddr == reserves.Token1 {
		// src = token0, dst = token1
		reserveIn = reserves.Reserve0
		reserveOut = reserves.Reserve1
	} else if srcAddr == reserves.Token1 && dstAddr == reserves.Token0 {
		// src = token1, dst = token0
		reserveIn = reserves.Reserve1
		reserveOut = reserves.Reserve0
	} else {
		return nil, models.NewAPIError(400, "Token mismatch", "Provided tokens don't match the pool tokens")
	}

	return &models.SwapCalculation{
		AmountIn:   amountIn,
		ReserveIn:  reserveIn,
		ReserveOut: reserveOut,
		TokenIn:    srcToken,
		TokenOut:   dstToken,
	}, nil
}
