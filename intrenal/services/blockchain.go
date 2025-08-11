package services

import (
	"context"
	"log"
	"math/big"
	"strings"
	"uniswap-est/intrenal/config"
	"uniswap-est/intrenal/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC20 ABI for decimals() and symbol() functions
const erc20ABI = `[
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [{"name": "", "type": "uint8"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol", 
		"outputs": [{"name": "", "type": "string"}],
		"type": "function"
	}
]`

// Uniswap V2 Pair ABI for getReserves() function
const pairABI = `[
	{
		"constant": true,
		"inputs": [],
		"name": "getReserves",
		"outputs": [
			{"name": "_reserve0", "type": "uint112"},
			{"name": "_reserve1", "type": "uint112"}, 
			{"name": "_blockTimestampLast", "type": "uint32"}
		],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "token0",
		"outputs": [{"name": "", "type": "address"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "token1", 
		"outputs": [{"name": "", "type": "address"}],
		"type": "function"
	}
]`

type BlockchainService struct {
	client   *ethclient.Client
	erc20ABI abi.ABI
	pairABI  abi.ABI
	config   *config.Config
}

// NewBlockchainService creates a new blockchain service
func NewBlockchainService(cfg *config.Config) (*BlockchainService, error) {
	// Connect to Ethereum node (Alchemy)
	client, err := ethclient.Dial(cfg.EthereumRPCURL)
	if err != nil {
		return nil, err
	}

	// Parse ABIs
	erc20Parsed, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, err
	}

	pairParsed, err := abi.JSON(strings.NewReader(pairABI))
	if err != nil {
		return nil, err
	}

	return &BlockchainService{
		client:   client,
		erc20ABI: erc20Parsed,
		pairABI:  pairParsed,
		config:   cfg,
	}, nil
}

// GetTokenInfo fetches token decimals and symbol from blockchain
func (bs *BlockchainService) GetTokenInfo(ctx context.Context, tokenAddress string) (*models.TokenInfo, error) {
	address := common.HexToAddress(tokenAddress)

	// Get decimals
	decimalsData, err := bs.erc20ABI.Pack("decimals")
	if err != nil {
		return nil, err
	}

	decimalsResult, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: decimalsData,
	}, nil)
	if err != nil {
		return nil, models.ErrBlockchainConnection
	}

	var decimals uint8
	err = bs.erc20ABI.UnpackIntoInterface(&decimals, "decimals", decimalsResult)
	if err != nil {
		return nil, err
	}

	// Get symbol
	symbolData, err := bs.erc20ABI.Pack("symbol")
	if err != nil {
		return nil, err
	}

	symbolResult, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: symbolData,
	}, nil)
	if err != nil {
		return nil, models.ErrBlockchainConnection
	}

	var symbol string
	err = bs.erc20ABI.UnpackIntoInterface(&symbol, "symbol", symbolResult)
	if err != nil {
		return nil, err
	}

	return &models.TokenInfo{
		Address:  tokenAddress,
		Decimals: decimals,
		Symbol:   symbol,
	}, nil
}

// GetPoolReserves fetches current reserves from Uniswap V2 pair
func (bs *BlockchainService) GetPoolReserves(ctx context.Context, poolAddress string) (*models.PoolReserves, error) {
	address := common.HexToAddress(poolAddress)

	// Debug logging
	log.Printf("Fetching reserves for pool: %s", poolAddress)
	log.Printf("Using RPC: %s", bs.config.EthereumRPCURL)

	// Get reserves
	reservesData, err := bs.pairABI.Pack("getReserves")
	if err != nil {
		log.Printf("Error packing getReserves: %v", err)
		return nil, err
	}

	log.Printf("Packed data length: %d bytes", len(reservesData))

	reservesResult, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: reservesData,
	}, nil)
	if err != nil {
		log.Printf("CallContract ERROR: %v", err)
		log.Printf("Pool address: %s", poolAddress)
		log.Printf("Call data: %x", reservesData)
		return nil, models.ErrPoolNotFound
	}

	log.Printf("Got reserves result: %x", reservesResult)

	var reserves struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}

	err = bs.pairABI.UnpackIntoInterface(&reserves, "getReserves", reservesResult)
	if err != nil {
		log.Printf("Error unpacking reserves: %v", err)
		return nil, err
	}

	// Get token0 address
	token0Data, err := bs.pairABI.Pack("token0")
	if err != nil {
		return nil, err
	}

	token0Result, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: token0Data,
	}, nil)
	if err != nil {
		return nil, models.ErrBlockchainConnection
	}

	var token0 common.Address
	err = bs.pairABI.UnpackIntoInterface(&token0, "token0", token0Result)
	if err != nil {
		return nil, err
	}

	// Get token1 address
	token1Data, err := bs.pairABI.Pack("token1")
	if err != nil {
		return nil, err
	}

	token1Result, err := bs.client.CallContract(ctx, ethereum.CallMsg{
		To:   &address,
		Data: token1Data,
	}, nil)
	if err != nil {
		return nil, models.ErrBlockchainConnection
	}

	var token1 common.Address
	err = bs.pairABI.UnpackIntoInterface(&token1, "token1", token1Result)
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully fetched pool data - Token0: %s, Token1: %s", token0.Hex(), token1.Hex())

	return &models.PoolReserves{
		Reserve0:  reserves.Reserve0,
		Reserve1:  reserves.Reserve1,
		Token0:    strings.ToLower(token0.Hex()),
		Token1:    strings.ToLower(token1.Hex()),
		BlockTime: reserves.BlockTimestampLast,
	}, nil
}

// Close closes the blockchain connection
func (bs *BlockchainService) Close() {
	bs.client.Close()
}