# DeFi Yield Aggregator

A full-stack application for tracking, analyzing, and discovering yield farming opportunities across multiple DeFi protocols in real-time.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![React](https://img.shields.io/badge/react-18+-61DAFB.svg)
![TypeScript](https://img.shields.io/badge/typescript-5.0+-3178C6.svg)

## Features

### Backend
- **Real-time Data Aggregation**: Fetches data from 2000+ pools across 20+ chains
- **Intelligent Opportunity Detection**: Identifies yield gaps, trending pools, and high-score opportunities
- **Risk-Adjusted Scoring**: Calculates scores based on APY, TVL, stability, and chain security
- **Time-Series Analytics**: Historical APY tracking with TimescaleDB
- **Fast Search**: Sub-second search across millions of records with ElasticSearch
- **Real-Time Updates**: WebSocket support for live data streaming
- **Production-Ready**: Docker, logging, error handling, and rate limiting

### Frontend
- **Modern React Dashboard**: Built with React 18, TypeScript, and Vite
- **500+ Mock Pools**: Comprehensive mock data across 20 chains and 40+ protocols
- **Interactive Charts**: APY history charts with Recharts and sparkline trends
- **External Protocol Links**: Direct deposit links to Aave, Compound, Uniswap, Curve, Lido, and 40+ protocols
- **Advanced Filtering**: Filter by chain, protocol, TVL, APY, stablecoins
- **Responsive Design**: Mobile-friendly dark theme with Tailwind CSS
- **Real-time Updates**: WebSocket integration for live pool updates

## Screenshots

The application includes:
- **Dashboard**: Overview of TVL, top pools, and active opportunities
- **Pools Explorer**: Table view with search, filters, and sorting
- **Pool Details**: Full APY history charts, token info, and quick deposit actions
- **Opportunities**: Yield gaps, trending pools, and arbitrage opportunities
- **Opportunity Details**: Comprehensive analysis with external links

## Tech Stack

### Backend
- **Go 1.21+** with Fiber framework
- **PostgreSQL 15** with TimescaleDB extension
- **ElasticSearch 8.x** for fast search and analytics
- **Redis 7.x** for caching and pub/sub
- **Docker Compose** for containerization

### Frontend
- **React 18** with TypeScript
- **Vite** for fast development and building
- **TanStack React Query v5** for data fetching
- **Tailwind CSS** for styling
- **Recharts** for data visualization
- **React Router v6** for navigation

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for frontend development)
- Go 1.21+ (for backend development)

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/defi-yield-aggregator.git
cd defi-yield-aggregator
```

### 2. Start Backend Services (Optional)

```bash
# Copy environment file
cp .env.example .env

# Start all services
docker-compose up -d
```

### 3. Start Frontend (with Mock Data)

The frontend works standalone with comprehensive mock data:

```bash
cd frontend
npm install
npm run dev
```

Open http://localhost:5173 in your browser.

### 4. Build for Production

```bash
# Frontend
cd frontend
npm run build

# Backend
docker build --target production-api -t defi-api:latest .
docker build --target production-worker -t defi-worker:latest .
```

## Project Structure

```
defi-yield-aggregator/
├── cmd/                          # Go entry points
│   ├── server/main.go           # HTTP server
│   └── worker/main.go           # Background worker
├── internal/                     # Go internal packages
│   ├── api/                     # HTTP handlers, middleware
│   ├── services/                # Business logic
│   ├── repository/              # Data access
│   └── models/                  # Data structures
├── frontend/                     # React frontend
│   ├── src/
│   │   ├── components/          # Reusable UI components
│   │   │   ├── APYChart.tsx    # APY history chart
│   │   │   ├── PoolCard.tsx    # Pool card with sparkline
│   │   │   ├── PoolTable.tsx   # Pools table view
│   │   │   ├── Sparkline.tsx   # Mini trend chart
│   │   │   ├── Icons.tsx       # Shared icons
│   │   │   └── ...
│   │   ├── pages/              # Page components
│   │   │   ├── HomePage.tsx
│   │   │   ├── PoolsPage.tsx
│   │   │   ├── PoolDetailsPage.tsx
│   │   │   ├── OpportunitiesPage.tsx
│   │   │   └── OpportunityDetailsPage.tsx
│   │   ├── services/           # API and mock data
│   │   ├── hooks/              # Custom React hooks
│   │   ├── utils/              # Utilities and helpers
│   │   │   ├── format.ts       # Number/date formatting
│   │   │   ├── links.ts        # Protocol URLs
│   │   │   └── constants.ts    # App constants
│   │   └── types/              # TypeScript types
│   └── package.json
├── migrations/                   # Database migrations
├── docker-compose.yml
└── README.md
```

## API Reference

### Health Check
```bash
GET /api/v1/health
```

### Pools
```bash
# List pools with filters
GET /api/v1/pools?chain=ethereum&minApy=5&sortBy=tvl&limit=20

# Get specific pool
GET /api/v1/pools/:id

# Get pool history
GET /api/v1/pools/:id/history?period=7d
```

### Opportunities
```bash
# List opportunities
GET /api/v1/opportunities?type=yield-gap&riskLevel=low

# Get trending pools
GET /api/v1/opportunities/trending
```

### Statistics
```bash
GET /api/v1/chains
GET /api/v1/protocols
GET /api/v1/stats
```

## Supported Protocols

The application supports 40+ DeFi protocols including:

| Category | Protocols |
|----------|-----------|
| Lending | Aave V3, Compound V3, Spark, Morpho, Radiant |
| DEXs | Uniswap V3, Curve, SushiSwap, Balancer, Velodrome |
| Liquid Staking | Lido, Rocket Pool, Frax ETH, Coinbase ETH |
| Yield Vaults | Yearn, Beefy, Convex, Pendle |
| Derivatives | GMX, dYdX, Synthetix, Gains Network |
| Bridges | Stargate, Hop Protocol, Across |

## Supported Chains

20 EVM-compatible chains:
- Ethereum, Arbitrum, Optimism, Base, Polygon
- Avalanche, BSC, Fantom, Gnosis, zkSync Era
- Polygon zkEVM, Linea, Scroll, Mantle, Blast
- Metis, Celo, Moonbeam, Aurora, Cronos

## External Links

Every pool and opportunity includes quick access to:
- **Protocol App**: Direct deposit link to the protocol
- **DefiLlama**: Compare yields and view analytics
- **Block Explorer**: View contracts on chain
- **DexScreener**: Check token trading data
- **CoinGecko**: Token information and prices
- **Zapper**: Portfolio tracking

## Configuration

### Frontend Environment Variables

Create `frontend/.env.local`:

```env
VITE_API_BASE=/api/v1
VITE_USE_MOCK_DATA=true
VITE_WS_BASE=ws://localhost:8080
```

### Backend Environment Variables

See `.env.example` for all options:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | 3000 |
| `POSTGRES_HOST` | PostgreSQL host | localhost |
| `REDIS_HOST` | Redis host | localhost |
| `ELASTICSEARCH_URL` | ElasticSearch URL | http://localhost:9200 |
| `DEFILLAMA_FETCH_INTERVAL` | Pool fetch interval | 3m |
| `MIN_TVL_THRESHOLD` | Minimum TVL to consider | 100000 |

## Development

### Frontend Development

```bash
cd frontend
npm install
npm run dev        # Start dev server
npm run build      # Production build
npm run lint       # Run ESLint
npm run preview    # Preview production build
```

### Backend Development

```bash
# Install dependencies
go mod download

# Run API server
go run cmd/server/main.go

# Run worker
go run cmd/worker/main.go

# Run with hot reload
air -c .air.toml

# Run tests
go test ./...
```

## Data Sources

### DeFiLlama
- Free API, no key required
- Updates every 3 minutes
- Provides 2000+ pools with APY, TVL data

### CoinGecko
- Free Demo plan available
- Updates every 10 minutes
- Token prices for calculations

## Opportunity Detection

### Yield Gap Arbitrage
Identifies the same asset with different APYs across protocols:
```
USDC on Aave (3.5%) vs Compound (4.2%) = +0.7% opportunity
```

### APY Trend Analysis
- Tracks 1h, 24h, 7d APY changes
- Alerts on significant APY increases (>50%)
- Identifies "hot pools" with growing yields

### Risk-Adjusted Scoring
```
score = (APY × 0.35) + (TVL × 0.25) + (Stability × 0.25) + (Trend × 0.15)
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [DeFiLlama](https://defillama.com) for yield data API
- [CoinGecko](https://coingecko.com) for price data
- [Fiber](https://gofiber.io) for Go web framework
- [TanStack Query](https://tanstack.com/query) for React data fetching
- [Tailwind CSS](https://tailwindcss.com) for styling
- [Recharts](https://recharts.org) for charts

## 👨‍💻 Author

Maksim Jatmanov - Backend Developer specializing in Go, Blockchain & AI