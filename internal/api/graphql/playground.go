package graphql

import (
	"github.com/gofiber/fiber/v2"
)

// PlaygroundHTML returns the GraphQL Playground HTML page
const PlaygroundHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
  <title>DeFi Yield Aggregator - GraphQL Playground</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
  <link rel="shortcut icon" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/favicon.png" />
  <script src="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>
<body>
  <div id="root">
    <style>
      body {
        background-color: rgb(23, 42, 58);
        font-family: Open Sans, sans-serif;
        height: 90vh;
      }
      #root {
        height: 100%;
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .loading {
        font-size: 32px;
        font-weight: 200;
        color: rgba(255, 255, 255, .6);
        margin-left: 28px;
      }
      img {
        width: 78px;
        height: 78px;
      }
      .title {
        font-weight: 400;
      }
    </style>
    <img src='https://cdn.jsdelivr.net/npm/graphql-playground-react/build/logo.png' alt=''>
    <div class="loading">
      Loading <span class="title">DeFi Yield Aggregator</span>
    </div>
  </div>
  <script>
    window.addEventListener('load', function (event) {
      GraphQLPlayground.init(document.getElementById('root'), {
        endpoint: '/graphql',
        settings: {
          'request.credentials': 'include',
          'editor.theme': 'dark',
          'editor.fontSize': 14,
          'editor.fontFamily': "'Source Code Pro', 'Consolas', 'Inconsolata', 'Droid Sans Mono', 'Monaco', monospace",
        },
        tabs: [
          {
            name: 'List Pools',
            endpoint: '/graphql',
            query: ` + "`" + `# List pools with filtering
query ListPools($filter: PoolFilter, $pagination: PaginationInput) {
  pools(filter: $filter, pagination: $pagination) {
    edges {
      node {
        id
        chain
        protocol
        symbol
        tvl
        apy
        score
        stablecoin
      }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
    totalCount
  }
}` + "`" + `,
            variables: JSON.stringify({
              filter: {
                chain: "ethereum",
                minApy: 5,
                sortBy: "TVL",
                sortOrder: "DESC"
              },
              pagination: {
                limit: 10,
                offset: 0
              }
            }, null, 2)
          },
          {
            name: 'Opportunities',
            endpoint: '/graphql',
            query: ` + "`" + `# List yield opportunities
query ListOpportunities($filter: OpportunityFilter) {
  opportunities(filter: $filter) {
    edges {
      node {
        id
        type
        title
        description
        asset
        chain
        apyDifference
        potentialProfit
        riskLevel
        score
        isActive
        detectedAt
      }
    }
    totalCount
  }
}` + "`" + `,
            variables: JSON.stringify({
              filter: {
                type: "YIELD_GAP",
                riskLevel: "LOW",
                activeOnly: true
              }
            }, null, 2)
          },
          {
            name: 'Trending Pools',
            endpoint: '/graphql',
            query: ` + "`" + `# Get trending pools with increasing APY
query TrendingPools($chain: String, $minGrowth: Float, $limit: Int) {
  trendingPools(chain: $chain, minGrowth: $minGrowth, limit: $limit) {
    pool {
      id
      chain
      protocol
      symbol
      apy
      tvl
    }
    apyGrowth1h
    apyGrowth24h
    apyGrowth7d
    trendScore
  }
}` + "`" + `,
            variables: JSON.stringify({
              minGrowth: 10,
              limit: 10
            }, null, 2)
          },
          {
            name: 'Platform Stats',
            endpoint: '/graphql',
            query: ` + "`" + `# Get platform-wide statistics
query PlatformStats {
  stats {
    totalPools
    totalTvl
    averageApy
    maxApy
    totalChains
    totalProtocols
    activeOpportunities
    lastUpdated
    tvlByChain {
      chain
      tvl
    }
    poolsByChain {
      chain
      count
    }
    apyDistribution {
      range0to1
      range1to5
      range5to10
      range10to25
      range25to50
      range50to100
      range100plus
    }
  }
}` + "`" + `
          },
          {
            name: 'Chains & Protocols',
            endpoint: '/graphql',
            query: ` + "`" + `# Get all chains and protocols
query ChainsAndProtocols {
  chains {
    name
    displayName
    poolCount
    totalTvl
    averageApy
    maxApy
  }
  protocols {
    edges {
      node {
        name
        displayName
        chains
        poolCount
        totalTvl
        averageApy
      }
    }
    totalCount
  }
}` + "`" + `
          },
          {
            name: 'Health Check',
            endpoint: '/graphql',
            query: ` + "`" + `# Check service health
query HealthCheck {
  health {
    status
    version
    uptime
    timestamp
    services {
      postgresql {
        status
        latency
      }
      redis {
        status
        latency
      }
      elasticsearch {
        status
        latency
      }
    }
  }
}` + "`" + `
          }
        ]
      })
    })
  </script>
</body>
</html>
`

// Playground returns the GraphQL Playground UI
func Playground(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(PlaygroundHTML)
}
