# DeFi Yield Aggregator

A production-ready, full-stack application for tracking, analyzing, and discovering yield farming opportunities across multiple DeFi protocols in real-time.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8.svg)
![React](https://img.shields.io/badge/react-18+-61DAFB.svg)
![TypeScript](https://img.shields.io/badge/typescript-5.0+-3178C6.svg)
![Docker](https://img.shields.io/badge/docker-ready-2496ED.svg)

## Features

### Backend (Go + Fiber)
- **Real-time Data Aggregation**: Fetches data from 2000+ pools across 30+ chains via DeFiLlama API
- **Intelligent Opportunity Detection**: Identifies yield gaps, trending pools, and high-score opportunities
- **Risk-Adjusted Scoring**: Calculates scores based on APY, TVL, stability, and chain security
- **Time-Series Analytics**: Historical APY tracking with TimescaleDB (PostgreSQL extension)
- **Fast Search**: Sub-100ms search across millions of records with ElasticSearch
- **Multi-Layer Caching**: Redis caching with intelligent cache key generation
- **Real-Time Updates**: WebSocket support for live pool/opportunity streaming
- **Production-Ready**: Docker, structured logging, error handling, rate limiting, and request timeouts

### Frontend (React 18 + TypeScript)
- **Modern Dashboard**: Built with React 18, TypeScript 5, and Vite 5
- **Type-Safe API Layer**: Full type transformation between API (strings) and frontend (numbers)
- **Interactive Charts**: APY history visualization with Recharts
- **Smart External Links**:
  - Protocol apps (40+ DeFi protocols)
  - Block explorers (30+ chain explorers)
  - CoinGecko token pages (symbol-to-ID mapping)
  - DexScreener charts
  - DeFiLlama analytics
- **Advanced Search**: Case-insensitive multi-field search (symbol, protocol, chain)
- **Advanced Filtering**: Filter by chain, protocol, TVL, APY, stablecoins with proper pagination
- **Infinite Scroll**: Load more with data accumulation for opportunities
- **Real-time Updates**: WebSocket integration with auto-reconnect
- **Responsive Design**: Mobile-friendly dark theme with Tailwind CSS

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Pages     │  │  Components │  │   Hooks     │              │
│  │  - Home     │  │  - PoolCard │  │  - useWS    │              │
│  │  - Pools    │  │  - Charts   │  │  - useQuery │              │
│  │  - Details  │  │  - Tables   │  │             │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│                           │                                      │
│                    API Service Layer                             │
│              (Type transformation layer)                         │
└──────────────────────────┬──────────────────────────────────────┘
                           │ HTTP/WebSocket
┌──────────────────────────┴──────────────────────────────────────┐
│                      API Server (Go/Fiber)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  Handlers   │  │ Middleware  │  │  WebSocket  │              │
│  │  - Pools    │  │  - CORS     │  │    Hub      │              │
│  │  - Opps     │  │  - Rate     │  │  (real-time)│              │
│  │  - Stats    │  │  - Timeout  │  │             │              │
│  └──────┬──────┘  └─────────────┘  └─────────────┘              │
│         │                                                        │
│  ┌──────┴─────────────────────────────────────────────┐         │
│  │              Service Layer                          │         │
│  │  - Opportunity Detection (yield gaps, trending)     │         │
│  │  - Scoring Engine (risk-adjusted scores)            │         │
│  │  - Pool Aggregation                                 │         │
│  └─────────────────────────┬───────────────────────────┘         │
└────────────────────────────┼────────────────────────────────────┘
                             │
┌────────────────────────────┼────────────────────────────────────┐
│                     Data Layer                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ PostgreSQL  │  │   Redis     │  │ElasticSearch│              │
│  │ +TimescaleDB│  │   Cache     │  │   Search    │              │
│  │ (primary)   │  │ (30-60s TTL)│  │ (fallback)  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                             │
┌────────────────────────────┼────────────────────────────────────┐
│                     Worker Service                               │
│  ┌─────────────────────────────────────────────┐                │
│  │  - Fetches pools from DeFiLlama (3m interval)│                │
│  │  - Indexes to ElasticSearch                  │                │
│  │  - Detects opportunities (5m interval)       │                │
│  │  - Broadcasts updates via WebSocket          │                │
│  └─────────────────────────────────────────────┘                │
└─────────────────────────────────────────────────────────────────┘
```

## Tech Stack

### Backend
| Component | Technology | Purpose |
|-----------|------------|---------|
| Framework | Go 1.23 + Fiber | High-performance HTTP server |
| Database | PostgreSQL 15 + TimescaleDB | Primary storage + time-series |
| Search | ElasticSearch 8.x | Fast full-text search |
| Cache | Redis 7.x | Response caching |
| Container | Docker + Docker Compose | Deployment |

### Frontend
| Component | Technology | Purpose |
|-----------|------------|---------|
| Framework | React 18 | UI components |
| Language | TypeScript 5 | Type safety |
| Build | Vite 5 | Fast development |
| Data | TanStack Query v5 | Server state management |
| Styling | Tailwind CSS 3 | Utility-first CSS |
| Charts | Recharts | Data visualization |
| Routing | React Router v6 | Navigation |

## Quick Start

### Prerequisites
- Docker and Docker Compose (for full stack)
- Node.js 18+ (for frontend only)
- Go 1.23+ (for backend development)

### Option 1: Full Stack with Docker

```bash
# Clone the repository
git clone https://github.com/yourusername/defi-yield-aggregator.git
cd defi-yield-aggregator

# Copy environment file
cp .env.example .env

# Start all services (API, Worker, PostgreSQL, Redis, ElasticSearch)
docker-compose up -d

# View logs
docker-compose logs -f api worker

# Access the application
# API: http://localhost:3000/api/v1/health
# Frontend: http://localhost:5173
```

### Option 2: Frontend Only (Mock Data)

```bash
cd frontend
npm install
npm run dev
```

Open http://localhost:5173 - works with comprehensive mock data.

### Option 3: Development Mode

```bash
# Terminal 1: Start infrastructure
docker-compose up -d postgres redis elasticsearch

# Terminal 2: Start API with hot reload
cd cmd/server && air

# Terminal 3: Start Worker
cd cmd/worker && air

# Terminal 4: Start Frontend
cd frontend && npm run dev
```

## API Reference

### Health & Stats
```bash
GET /api/v1/health              # Service health check
GET /api/v1/stats               # Aggregated statistics
GET /api/v1/chains              # List of supported chains
GET /api/v1/protocols           # List of protocols
```

### Pools
```bash
# List pools with filters
GET /api/v1/pools
  ?chain=ethereum               # Filter by chain (case-insensitive)
  &protocol=aave-v3             # Filter by protocol
  &symbol=ETH                   # Filter by symbol (partial match)
  &search=USDC                  # Search across all fields
  &minApy=5                     # Minimum APY
  &maxApy=100                   # Maximum APY
  &minTvl=1000000              # Minimum TVL
  &minScore=50                  # Minimum score
  &stablecoin=true             # Stablecoin pools only
  &sortBy=apy|tvl|score        # Sort field (default: tvl)
  &sortOrder=asc|desc          # Sort order (default: desc)
  &limit=50                     # Results per page (max: 100)
  &offset=0                     # Pagination offset

# Get specific pool
GET /api/v1/pools/:id

# Get pool APY history
GET /api/v1/pools/:id/history
  ?period=1h|24h|7d|30d        # Time period (default: 24h)
```

### Opportunities
```bash
# List opportunities
GET /api/v1/opportunities
  ?type=yield-gap|trending|high-score
  &riskLevel=low|medium|high
  &chain=ethereum
  &asset=USDC
  &minProfit=1
  &minScore=50
  &activeOnly=true             # Active opportunities only
  &sortBy=score|profit|apy     # Sort field
  &limit=50
  &offset=0

# Get trending pools
GET /api/v1/opportunities/trending
  ?chain=ethereum
  &minGrowth=10                # Minimum APY growth %
  &limit=20
```

### WebSocket
```javascript
// Connect to pools stream
ws://localhost:3000/ws/pools

// Connect to opportunities stream
ws://localhost:3000/ws/opportunities

// Message types received:
// - pool_update: Real-time pool data changes
// - opportunity_alert: New opportunity detected
// - ping/pong: Keep-alive
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| **Server** |||
| `SERVER_PORT` | API server port | 3000 |
| `SERVER_READ_TIMEOUT` | Request read timeout | 30s |
| `APP_ENV` | Environment (development/production) | development |
| **Database** |||
| `POSTGRES_HOST` | PostgreSQL host | localhost |
| `POSTGRES_PORT` | PostgreSQL port | 5432 |
| `POSTGRES_USER` | Database user | defi |
| `POSTGRES_PASSWORD` | Database password | ⚠️ Change in production |
| `POSTGRES_DB` | Database name | defi_aggregator |
| `POSTGRES_MAX_CONNECTIONS` | Connection pool size | 25 |
| **Redis** |||
| `REDIS_HOST` | Redis host | localhost |
| `REDIS_PORT` | Redis port | 6379 |
| `REDIS_POOL_SIZE` | Connection pool size | 10 |
| **ElasticSearch** |||
| `ELASTICSEARCH_URL` | ElasticSearch URL | http://localhost:9200 |
| **Data Fetching** |||
| `DEFILLAMA_FETCH_INTERVAL` | Pool fetch interval | 3m |
| `OPPORTUNITY_DETECT_INTERVAL` | Opportunity detection interval | 5m |
| `MIN_TVL_THRESHOLD` | Minimum TVL to consider | 100000 |
| `MIN_APY_THRESHOLD` | Minimum APY to consider | 0.1 |
| `YIELD_GAP_MIN_PROFIT` | Min profit for yield gap alerts | 0.5 |
| **Rate Limiting** |||
| `RATE_LIMIT_REQUESTS` | Requests per window | 100 |
| `RATE_LIMIT_WINDOW` | Rate limit window | 1m |
| **CORS** |||
| `CORS_ALLOWED_ORIGINS` | Allowed origins | * (⚠️ Restrict in production) |

### Frontend Configuration

Create `frontend/.env.local`:
```env
VITE_API_BASE=/api/v1
VITE_USE_MOCK_DATA=false
```

## Supported Protocols (40+)

| Category | Protocols |
|----------|-----------|
| **Lending** | Aave V3, Compound V3, Spark, Morpho, Radiant, Venus, BenQi |
| **DEXs** | Uniswap V3/V2, Curve, SushiSwap, Balancer, Velodrome, Aerodrome, PancakeSwap |
| **Liquid Staking** | Lido, Rocket Pool, Frax ETH, Coinbase ETH, EigenLayer, EtherFi, Renzo |
| **Yield Vaults** | Yearn, Beefy, Convex, Pendle, Sommelier |
| **Derivatives** | GMX, Gains Network |
| **Bridges** | Stargate, Hop Protocol, Across |
| **Institutional** | Maple, Goldfinch |

## Supported Chains (30+)

| Category | Chains |
|----------|--------|
| **L1** | Ethereum, Avalanche, BSC, Fantom, Solana, Aptos, Sui |
| **L2 Optimistic** | Arbitrum, Optimism, Base, Mantle, Blast, Mode |
| **L2 ZK** | zkSync Era, Polygon zkEVM, Linea, Scroll, StarkNet |
| **Alt L1** | Polygon, Gnosis, Celo, Moonbeam, Kava, Aurora, Cronos |
| **New** | Berachain, Sonic, Hyperliquid, Monad |

## Opportunity Detection

### Yield Gap Arbitrage
Identifies the same asset with different APYs across protocols:
```
USDC on Aave V3 (3.5% APY) vs Compound V3 (4.2% APY)
→ +0.7% yield gap opportunity
→ Estimated profit: $7,000/year on $1M position
```

### Trending Pools
Detects pools with rapidly increasing APY:
```
Pool: WETH on Aerodrome
→ APY increased 150% in 24h
→ Current APY: 25.5%
→ TVL: $50M
```

### Risk-Adjusted Scoring
```
Score = (APY × 0.35) + (TVL × 0.25) + (Stability × 0.25) + (Trend × 0.15)

Where:
- APY: Normalized yield percentage
- TVL: Liquidity depth indicator
- Stability: 30-day APY variance
- Trend: 7-day momentum
```

## Performance Optimizations

### Backend
- **Connection Pooling**: PostgreSQL (25 connections), Redis (10 connections)
- **Request Timeouts**: 30-second context timeout on all database operations
- **Multi-Layer Caching**: Redis cache with comprehensive cache keys
- **ElasticSearch Fallback**: Automatic fallback to PostgreSQL if ES returns no results
- **WebSocket Optimization**: Dead client cleanup, race condition fixes

### Frontend
- **Type Transformation**: API returns strings, frontend converts to numbers
- **Data Accumulation**: Infinite scroll without losing previous data
- **Optimistic Updates**: React Query with placeholder data
- **Debounced Search**: 300ms debounce on search input

## Project Structure

```
defi-yield-aggregator/
├── cmd/
│   ├── server/main.go          # API server entry point
│   └── worker/main.go          # Background worker entry point
├── internal/
│   ├── api/
│   │   ├── handlers/           # HTTP handlers with validation
│   │   ├── middleware/         # CORS, rate limiting, logging
│   │   └── websocket/          # WebSocket hub and clients
│   ├── config/                 # Configuration management
│   ├── models/                 # Data structures
│   ├── repository/
│   │   ├── postgres/           # PostgreSQL + TimescaleDB
│   │   ├── redis/              # Redis caching
│   │   └── elasticsearch/      # ElasticSearch search
│   └── services/
│       ├── defillama/          # DeFiLlama API client
│       ├── opportunity/        # Opportunity detection
│       └── scoring/            # Risk scoring engine
├── frontend/
│   ├── src/
│   │   ├── components/         # Reusable UI components
│   │   ├── pages/              # Page components
│   │   ├── services/           # API client with transformations
│   │   ├── hooks/              # Custom React hooks
│   │   ├── utils/              # Helpers (links, formatting)
│   │   └── types/              # TypeScript definitions
│   └── package.json
├── migrations/                 # Database migrations
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Development

### Running Tests
```bash
# Backend tests
go test ./... -v

# Frontend tests
cd frontend && npm test
```

### Building for Production
```bash
# Backend
docker build --target production-api -t defi-api:latest .
docker build --target production-worker -t defi-worker:latest .

# Frontend
cd frontend && npm run build
```

### Code Quality
```bash
# Go
go fmt ./...
go vet ./...
golangci-lint run

# TypeScript
cd frontend
npm run lint
npm run type-check
```

## Security Considerations

⚠️ **Before deploying to production:**

1. **Change default credentials** in `.env`
2. **Restrict CORS origins** (`CORS_ALLOWED_ORIGINS`)
3. **Enable TLS/SSL** for all connections
4. **Set up proper authentication** for sensitive endpoints
5. **Configure rate limiting** appropriately
6. **Use secrets management** (Vault, AWS Secrets, etc.)

## Troubleshooting

### Common Issues

**Search not working:**
- Check ElasticSearch is running: `curl localhost:9200/_cluster/health`
- Verify index exists: `curl localhost:9200/defi_pools/_count`
- Check API logs for errors

**Opportunities not loading:**
- Ensure worker is running and detecting opportunities
- Check PostgreSQL has opportunity data
- Verify cache is not stale

**WebSocket disconnects:**
- Check for network issues
- Verify `WS_PING_INTERVAL` and `WS_PONG_TIMEOUT` settings
- Check browser console for errors

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [DeFiLlama](https://defillama.com) - Yield data API
- [CoinGecko](https://coingecko.com) - Token data
- [Fiber](https://gofiber.io) - Go web framework
- [TanStack Query](https://tanstack.com/query) - React data fetching
- [Tailwind CSS](https://tailwindcss.com) - Styling
- [Recharts](https://recharts.org) - Charts

---

**Author**: Maksim Jatmanov - Backend Developer specializing in Go, Blockchain & AI
