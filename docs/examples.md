# API Examples

## GraphQL API

Access the GraphQL Playground at: http://localhost:3000/graphql

### Query Pools with GraphQL

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { pools(filter: { chain: \"ethereum\", minApy: 5 }, pagination: { limit: 10 }) { edges { node { id chain protocol symbol tvl apy score } } totalCount } }"
  }' | jq
```

### Query Opportunities

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { opportunities(filter: { type: YIELD_GAP, activeOnly: true }) { edges { node { id type title apyDifference potentialProfit riskLevel } } totalCount } }"
  }' | jq
```

### Query Platform Stats

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { stats { totalPools totalTvl averageApy totalChains totalProtocols activeOpportunities } chains { name poolCount totalTvl averageApy } }"
  }' | jq
```

### Query Trending Pools

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { trendingPools(minGrowth: 10, limit: 5) { pool { id symbol apy } apyGrowth24h trendScore } }"
  }' | jq
```

### Health Check via GraphQL

```bash
curl -X POST http://localhost:3000/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { health { status version uptime services { postgresql { status latency } redis { status latency } elasticsearch { status latency } } } }"
  }' | jq
```

---

## REST API

## Health Check

```bash
# Check service health
curl http://localhost:3000/api/v1/health | jq
```

Response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "postgresql": {
      "status": "up",
      "latency": "2.5ms"
    },
    "redis": {
      "status": "up",
      "latency": "0.5ms"
    },
    "elasticsearch": {
      "status": "up",
      "latency": "5ms"
    }
  }
}
```

## List Pools

```bash
# Get top 10 pools by TVL on Ethereum
curl "http://localhost:3000/api/v1/pools?chain=ethereum&sortBy=tvl&limit=10" | jq

# Get stablecoin pools with APY > 5%
curl "http://localhost:3000/api/v1/pools?stablecoin=true&minApy=5&sortBy=apy&sortOrder=desc" | jq

# Get high-score pools (score > 70)
curl "http://localhost:3000/api/v1/pools?minScore=70&sortBy=score&sortOrder=desc&limit=20" | jq

# Search for USDC pools
curl "http://localhost:3000/api/v1/pools?symbol=USDC" | jq

# Get Aave V3 pools
curl "http://localhost:3000/api/v1/pools?protocol=aave-v3" | jq
```

## Get Single Pool

```bash
# Get pool details
curl "http://localhost:3000/api/v1/pools/aave-v3-ethereum-usdc" | jq
```

Response:
```json
{
  "id": "aave-v3-ethereum-usdc",
  "chain": "ethereum",
  "protocol": "aave-v3",
  "symbol": "USDC",
  "tvl": 500000000,
  "apy": 3.5,
  "apyBase": 2.5,
  "apyReward": 1.0,
  "score": 85.5,
  "apyChange1h": 0.02,
  "apyChange24h": 0.15,
  "apyChange7d": -0.3,
  "stablecoin": true,
  "rewardTokens": ["AAVE"],
  "underlyingTokens": ["USDC"],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

## Pool History

```bash
# Get 24-hour APY history
curl "http://localhost:3000/api/v1/pools/aave-v3-ethereum-usdc/history?period=24h" | jq

# Get 7-day history
curl "http://localhost:3000/api/v1/pools/aave-v3-ethereum-usdc/history?period=7d" | jq
```

Response:
```json
{
  "poolId": "aave-v3-ethereum-usdc",
  "period": "24h",
  "dataPoints": [
    {
      "timestamp": "2024-01-14T10:00:00Z",
      "apy": 3.45,
      "tvl": 498000000
    },
    {
      "timestamp": "2024-01-14T10:05:00Z",
      "apy": 3.47,
      "tvl": 499000000
    }
  ]
}
```

## List Opportunities

```bash
# Get all active opportunities
curl "http://localhost:3000/api/v1/opportunities" | jq

# Get yield-gap opportunities only
curl "http://localhost:3000/api/v1/opportunities?type=yield-gap" | jq

# Get low-risk opportunities
curl "http://localhost:3000/api/v1/opportunities?riskLevel=low&sortBy=profit&sortOrder=desc" | jq

# Get opportunities for USDC
curl "http://localhost:3000/api/v1/opportunities?asset=USDC" | jq
```

Response:
```json
{
  "data": [
    {
      "id": "opp-abc123",
      "type": "yield-gap",
      "title": "USDC Yield Gap: 0.70% difference",
      "description": "Move USDC from Aave (3.5%) to Compound (4.2%)",
      "sourcePoolId": "aave-v3-ethereum-usdc",
      "targetPoolId": "compound-v3-ethereum-usdc",
      "asset": "USDC",
      "chain": "ethereum",
      "apyDifference": 0.7,
      "currentApy": 4.2,
      "potentialProfit": 175.0,
      "riskLevel": "low",
      "score": 82.5,
      "isActive": true,
      "detectedAt": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 45,
  "limit": 50,
  "offset": 0,
  "hasMore": false
}
```

## Trending Pools

```bash
# Get trending pools (APY growth > 10%)
curl "http://localhost:3000/api/v1/opportunities/trending" | jq

# Get trending pools on specific chain
curl "http://localhost:3000/api/v1/opportunities/trending?chain=arbitrum&minGrowth=20" | jq
```

## List Chains

```bash
curl "http://localhost:3000/api/v1/chains" | jq
```

Response:
```json
{
  "data": [
    {
      "name": "ethereum",
      "displayName": "Ethereum",
      "poolCount": 500,
      "totalTvl": 25000000000,
      "averageApy": 4.5,
      "maxApy": 25.0,
      "topProtocols": ["aave-v3", "compound-v3", "curve"]
    },
    {
      "name": "arbitrum",
      "displayName": "Arbitrum",
      "poolCount": 200,
      "totalTvl": 5000000000,
      "averageApy": 6.2,
      "maxApy": 50.0,
      "topProtocols": ["gmx", "aave-v3", "radiant"]
    }
  ],
  "total": 25
}
```

## List Protocols

```bash
# Get all protocols sorted by TVL
curl "http://localhost:3000/api/v1/protocols?sortBy=tvl" | jq

# Get protocols on Ethereum
curl "http://localhost:3000/api/v1/protocols?chain=ethereum" | jq
```

## Platform Statistics

```bash
curl "http://localhost:3000/api/v1/stats" | jq
```

Response:
```json
{
  "totalPools": 2500,
  "totalTvl": 50000000000,
  "averageApy": 5.5,
  "medianApy": 3.2,
  "maxApy": 500.0,
  "totalChains": 25,
  "totalProtocols": 150,
  "activeOpportunities": 45,
  "lastUpdated": "2024-01-15T10:30:00Z",
  "tvlByChain": {
    "ethereum": 25000000000,
    "bsc": 8000000000,
    "polygon": 5000000000
  },
  "poolsByChain": {
    "ethereum": 500,
    "bsc": 400,
    "polygon": 300
  },
  "apyDistribution": {
    "range0to1": 150,
    "range1to5": 800,
    "range5to10": 500,
    "range10to25": 400,
    "range25to50": 300,
    "range50to100": 200,
    "range100plus": 150
  }
}
```

## Metrics

```bash
# JSON metrics
curl "http://localhost:3000/api/v1/metrics" | jq

# Prometheus format
curl "http://localhost:3000/metrics"
```

## WebSocket Examples

### Connect to Pool Updates

```javascript
// JavaScript WebSocket client example
const ws = new WebSocket('ws://localhost:3000/ws/pools');

ws.onopen = () => {
  console.log('Connected to pool updates');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Pool update:', message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected from pool updates');
};
```

### Connect to Opportunity Alerts

```javascript
const ws = new WebSocket('ws://localhost:3000/ws/opportunities');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  if (message.type === 'opportunity_alert') {
    console.log('New opportunity detected!', message.data);
  }
};
```

### Using wscat (CLI tool)

```bash
# Install wscat
npm install -g wscat

# Connect to pool updates
wscat -c ws://localhost:3000/ws/pools

# Connect to opportunity alerts
wscat -c ws://localhost:3000/ws/opportunities
```

## Error Responses

### Validation Error (422)

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "One or more fields failed validation",
    "errors": [
      {
        "field": "minApy",
        "message": "must be a valid number"
      },
      {
        "field": "sortBy",
        "message": "invalid sort field"
      }
    ]
  }
}
```

### Not Found Error (404)

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Resource not found",
    "details": "Pool 'invalid-pool-id' not found"
  }
}
```

### Rate Limit Error (429)

```json
{
  "error": {
    "code": "RATE_LIMITED",
    "message": "Rate limit exceeded. Please try again later."
  }
}
```
