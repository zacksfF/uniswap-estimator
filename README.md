# 🦄 Uniswap V2 Estimator API

> fast REST API for estimating Uniswap V2 token swaps with real-time blockchain data integration.

[![Go](https://img.shields.io/badge/Go-1.24.5-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=flat)](https://gofiber.io/)
[![Ethereum](https://img.shields.io/badge/Ethereum-Compatible-627EEA?style=flat&logo=ethereum)](https://ethereum.org/)

## Features

- Single `/estimate` endpoint for swap calculations
- Real-time blockchain state fetching from Ethereum mainnet
- Custom Uniswap V2 math implementation with 0.3% fee calculation
- Performance-optimized with minimal memory allocations
- Built with Go and Fiber framework
- Comprehensive input validation and error handling

## Quick Start

1. **Clone and setup**
   ```bash
   git clone https://github.com/zacksfF/uniswap-estimator.git
   cd uniswap-estimator
   cp .env.example .env
   ```

2. **Install dependencies and run**
   ```bash
   make run
   ```

## Configure Environment

Edit `.env` file with your Ethereum RPC endpoint:

```env
ETHEREUM_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/your_key
HOST=localhost
PORT=1337
ENV=development
REQUEST_TIMEOUT=10
MAX_CONNECTIONS=100
```

Get RPC URL from:
- [Alchemy](https://alchemy.com)

## Commands / Development

```bash
make run        # Start development server
make build      # Build production binary
make test       # Run test suite
make benchmark  # Run performance benchmarks
make clean      # Clean build artifacts
make help       # Show all commands
```

**Performance benchmarks:**
```
BenchmarkCalculateAmountOut-8    7968488    135.0 ns/op    248 B/op    6 allocs/op
```

## Test the API

**Health check:**
```bash
curl http://localhost:1337/health
```

**Estimate endpoint:**
```bash
curl "http://localhost:1337/estimate?pool=0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852&src=0xdAC17F958D2ee523a2206206994597C13D831ec7&dst=0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2&src_amount=10000000"
```

**Expected response:**
```json
{
  "dst_amount": "238316708782106591"
}
```

**API parameters:**
- `pool` - Uniswap V2 pair contract address (42 chars)
- `src` - Source token address (42 chars)
- `dst` - Destination token address (42 chars)
- `src_amount` - Input amount as integer string

## Project Architecture

```
├── cmd/main.go                 # Application entry point
├── internal/
│   ├── config/                 # Environment configuration
│   ├── handlers/               # HTTP request handlers
│   ├── services/               # Business logic layer
│   ├── models/                 # Data structures and errors
│   └── utils/                  # Math and validation utilities
├── test/                       # Tests and benchmarks
├── .env                        # Environment variables
└── Makefile                    # Build automation
```

**Architecture principles:**
- Clean separation between blockchain interaction and mathematical calculation
- Fiber framework for high-performance HTTP handling
- Comprehensive error handling and input validation
- Production-ready with health monitoring

## Technical Implementation

**Blockchain Integration:**
- Fetches current pool reserves via Ethereum RPC calls
- Retrieves token metadata (decimals, symbols) from ERC20 contracts
- Uses raw contract calls without external Uniswap libraries
- Real-time data ensures accurate calculations

**Performance Optimization:**
- 135 nanoseconds per calculation
- 7.4 million calculations per second throughput
- 248 bytes allocated per operation
- Only 6 memory allocations per calculation

## Mathematical Engine

**Uniswap V2 Formula Implementation:**
- Constant product formula: `x * y = k`
- Trading fee calculation: `(amountIn * 997 * reserveOut) / (reserveIn * 1000 + amountIn * 997)`
- 0.3% fee (997/1000 factor) applied to input amount
- High precision using Go's `math/big` package
- Optimized for minimal memory allocations

**Token Handling:**
- Automatic decimal conversion for different token standards
- Proper token ordering based on Uniswap V2 pair structure
- Comprehensive amount validation

## Requirements Compliance

- Single `/estimate` endpoint with specified parameters
- No on-chain function calls for amount calculations
- No external Uniswap SDK or calculation libraries
- Real-time blockchain state fetching only
- Backend mathematical calculations with performance optimization
- Go implementation with web framework (Fiber)