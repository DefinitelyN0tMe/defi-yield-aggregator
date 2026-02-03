import type { Pool, Opportunity, PlatformStats, Chain, TrendingPool } from '../types';

// Real DeFi protocols and their categories
const protocols = [
  { name: 'aave-v3', displayName: 'Aave V3', category: 'lending' },
  { name: 'aave-v2', displayName: 'Aave V2', category: 'lending' },
  { name: 'compound-v3', displayName: 'Compound V3', category: 'lending' },
  { name: 'compound-v2', displayName: 'Compound V2', category: 'lending' },
  { name: 'lido', displayName: 'Lido', category: 'liquid-staking' },
  { name: 'rocket-pool', displayName: 'Rocket Pool', category: 'liquid-staking' },
  { name: 'frax-ether', displayName: 'Frax Ether', category: 'liquid-staking' },
  { name: 'curve', displayName: 'Curve', category: 'dex' },
  { name: 'convex', displayName: 'Convex', category: 'yield' },
  { name: 'uniswap-v3', displayName: 'Uniswap V3', category: 'dex' },
  { name: 'uniswap-v2', displayName: 'Uniswap V2', category: 'dex' },
  { name: 'sushiswap', displayName: 'SushiSwap', category: 'dex' },
  { name: 'balancer-v2', displayName: 'Balancer V2', category: 'dex' },
  { name: 'gmx', displayName: 'GMX', category: 'derivatives' },
  { name: 'gains-network', displayName: 'Gains Network', category: 'derivatives' },
  { name: 'pendle', displayName: 'Pendle', category: 'yield' },
  { name: 'yearn', displayName: 'Yearn', category: 'yield' },
  { name: 'beefy', displayName: 'Beefy', category: 'yield' },
  { name: 'velodrome', displayName: 'Velodrome', category: 'dex' },
  { name: 'aerodrome', displayName: 'Aerodrome', category: 'dex' },
  { name: 'camelot', displayName: 'Camelot', category: 'dex' },
  { name: 'trader-joe', displayName: 'Trader Joe', category: 'dex' },
  { name: 'pancakeswap', displayName: 'PancakeSwap', category: 'dex' },
  { name: 'quickswap', displayName: 'QuickSwap', category: 'dex' },
  { name: 'radiant', displayName: 'Radiant', category: 'lending' },
  { name: 'morpho', displayName: 'Morpho', category: 'lending' },
  { name: 'spark', displayName: 'Spark', category: 'lending' },
  { name: 'venus', displayName: 'Venus', category: 'lending' },
  { name: 'benqi', displayName: 'Benqi', category: 'lending' },
  { name: 'stargate', displayName: 'Stargate', category: 'bridge' },
  { name: 'hop-protocol', displayName: 'Hop Protocol', category: 'bridge' },
  { name: 'across', displayName: 'Across', category: 'bridge' },
  { name: 'eigenlayer', displayName: 'EigenLayer', category: 'restaking' },
  { name: 'ether-fi', displayName: 'Ether.fi', category: 'liquid-staking' },
  { name: 'renzo', displayName: 'Renzo', category: 'restaking' },
  { name: 'kelp-dao', displayName: 'Kelp DAO', category: 'restaking' },
  { name: 'sommelier', displayName: 'Sommelier', category: 'yield' },
  { name: 'instadapp', displayName: 'Instadapp', category: 'yield' },
  { name: 'maple', displayName: 'Maple', category: 'lending' },
  { name: 'goldfinch', displayName: 'Goldfinch', category: 'lending' },
];

// Chains with their characteristics
const chains = [
  { name: 'ethereum', displayName: 'Ethereum', tvlMultiplier: 1, apyMultiplier: 0.8 },
  { name: 'arbitrum', displayName: 'Arbitrum', tvlMultiplier: 0.3, apyMultiplier: 1.2 },
  { name: 'optimism', displayName: 'Optimism', tvlMultiplier: 0.15, apyMultiplier: 1.3 },
  { name: 'polygon', displayName: 'Polygon', tvlMultiplier: 0.12, apyMultiplier: 1.1 },
  { name: 'base', displayName: 'Base', tvlMultiplier: 0.1, apyMultiplier: 1.4 },
  { name: 'bsc', displayName: 'BNB Chain', tvlMultiplier: 0.08, apyMultiplier: 1.5 },
  { name: 'avalanche', displayName: 'Avalanche', tvlMultiplier: 0.06, apyMultiplier: 1.2 },
  { name: 'fantom', displayName: 'Fantom', tvlMultiplier: 0.02, apyMultiplier: 2.0 },
  { name: 'gnosis', displayName: 'Gnosis', tvlMultiplier: 0.015, apyMultiplier: 1.1 },
  { name: 'zkSync', displayName: 'zkSync Era', tvlMultiplier: 0.05, apyMultiplier: 1.6 },
  { name: 'linea', displayName: 'Linea', tvlMultiplier: 0.03, apyMultiplier: 1.8 },
  { name: 'scroll', displayName: 'Scroll', tvlMultiplier: 0.02, apyMultiplier: 2.0 },
  { name: 'mantle', displayName: 'Mantle', tvlMultiplier: 0.025, apyMultiplier: 1.7 },
  { name: 'manta', displayName: 'Manta', tvlMultiplier: 0.015, apyMultiplier: 2.2 },
  { name: 'blast', displayName: 'Blast', tvlMultiplier: 0.04, apyMultiplier: 2.5 },
  { name: 'mode', displayName: 'Mode', tvlMultiplier: 0.01, apyMultiplier: 3.0 },
  { name: 'metis', displayName: 'Metis', tvlMultiplier: 0.01, apyMultiplier: 1.8 },
  { name: 'celo', displayName: 'Celo', tvlMultiplier: 0.008, apyMultiplier: 1.5 },
  { name: 'moonbeam', displayName: 'Moonbeam', tvlMultiplier: 0.005, apyMultiplier: 1.6 },
  { name: 'kava', displayName: 'Kava', tvlMultiplier: 0.007, apyMultiplier: 1.9 },
];

// Token pairs and singles
const tokens = {
  stablecoins: ['USDC', 'USDT', 'DAI', 'FRAX', 'LUSD', 'crvUSD', 'GHO', 'PYUSD', 'TUSD', 'BUSD', 'USDP', 'sUSD', 'MIM', 'DOLA', 'alUSD'],
  majors: ['ETH', 'WETH', 'stETH', 'wstETH', 'rETH', 'cbETH', 'frxETH', 'sfrxETH', 'WBTC', 'BTC', 'tBTC'],
  altcoins: ['LINK', 'UNI', 'AAVE', 'CRV', 'CVX', 'LDO', 'RPL', 'GMX', 'ARB', 'OP', 'MATIC', 'SNX', 'COMP', 'MKR', 'BAL', 'SUSHI', 'YFI', 'FXS', 'SPELL', 'ALCX'],
  lsts: ['stETH', 'wstETH', 'rETH', 'cbETH', 'frxETH', 'sfrxETH', 'ankrETH', 'swETH', 'ETHx', 'mETH', 'eETH', 'weETH', 'ezETH', 'pufETH', 'rsETH'],
};

const rewardTokens: Record<string, string[]> = {
  'aave-v3': ['AAVE'],
  'aave-v2': ['AAVE', 'stkAAVE'],
  'compound-v3': ['COMP'],
  'compound-v2': ['COMP'],
  'curve': ['CRV'],
  'convex': ['CVX', 'CRV'],
  'uniswap-v3': [],
  'uniswap-v2': [],
  'sushiswap': ['SUSHI'],
  'balancer-v2': ['BAL'],
  'gmx': ['esGMX', 'ETH'],
  'pendle': ['PENDLE'],
  'yearn': ['YFI'],
  'beefy': ['BIFI'],
  'velodrome': ['VELO'],
  'aerodrome': ['AERO'],
  'camelot': ['GRAIL', 'xGRAIL'],
  'trader-joe': ['JOE'],
  'pancakeswap': ['CAKE'],
  'quickswap': ['QUICK'],
  'radiant': ['RDNT'],
  'stargate': ['STG'],
  'lido': [],
  'rocket-pool': [],
};

// Generate a random number in range with optional decimal places
function random(min: number, max: number, decimals = 2): number {
  return Number((Math.random() * (max - min) + min).toFixed(decimals));
}

// Generate unique pool ID
function generatePoolId(protocol: string, chain: string, symbol: string): string {
  return `${protocol}-${chain}-${symbol.toLowerCase().replace(/[^a-z0-9]/g, '-')}`;
}

// Generate pools
function generatePools(): Pool[] {
  const pools: Pool[] = [];
  const now = new Date();

  // Generate lending pools
  const lendingProtocols = protocols.filter(p => p.category === 'lending');
  for (const protocol of lendingProtocols) {
    for (const chain of chains) {
      // Skip some chain/protocol combinations randomly
      if (Math.random() > 0.6) continue;

      // Stablecoin lending
      for (const token of tokens.stablecoins.slice(0, Math.floor(Math.random() * 8) + 3)) {
        if (Math.random() > 0.7) continue;
        const baseApy = random(2, 12) * chain.apyMultiplier;
        const rewardApy = Math.random() > 0.5 ? random(0.5, 5) : 0;
        const tvl = random(1000000, 500000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, token),
          chain: chain.name,
          protocol: protocol.name,
          symbol: token,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || [],
          underlyingTokens: [token],
          poolMeta: `${protocol.displayName} ${token} Lending`,
          il7d: 0,
          apyMean30d: baseApy + rewardApy + random(-1, 1),
          volumeUsd1d: tvl * random(0.01, 0.1),
          volumeUsd7d: tvl * random(0.05, 0.5),
          score: random(60, 95, 0),
          apyChange1h: random(-0.5, 0.5),
          apyChange24h: random(-2, 2),
          apyChange7d: random(-5, 5),
          stablecoin: true,
          exposure: 'single',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }

      // Major token lending (ETH, WBTC, etc.)
      for (const token of tokens.majors.slice(0, Math.floor(Math.random() * 5) + 2)) {
        if (Math.random() > 0.6) continue;
        const baseApy = random(1, 8) * chain.apyMultiplier;
        const rewardApy = Math.random() > 0.5 ? random(0.3, 3) : 0;
        const tvl = random(5000000, 2000000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, token),
          chain: chain.name,
          protocol: protocol.name,
          symbol: token,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || [],
          underlyingTokens: [token],
          poolMeta: `${protocol.displayName} ${token} Lending`,
          il7d: 0,
          apyMean30d: baseApy + rewardApy + random(-0.5, 0.5),
          volumeUsd1d: tvl * random(0.01, 0.08),
          volumeUsd7d: tvl * random(0.05, 0.4),
          score: random(65, 98, 0),
          apyChange1h: random(-0.3, 0.3),
          apyChange24h: random(-1.5, 1.5),
          apyChange7d: random(-4, 4),
          stablecoin: false,
          exposure: 'single',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }
    }
  }

  // Generate DEX LP pools
  const dexProtocols = protocols.filter(p => p.category === 'dex');
  for (const protocol of dexProtocols) {
    for (const chain of chains) {
      if (Math.random() > 0.5) continue;

      // Stable pairs
      const stablePairs = [
        ['USDC', 'USDT'], ['DAI', 'USDC'], ['FRAX', 'USDC'], ['LUSD', 'USDC'],
        ['crvUSD', 'USDC'], ['GHO', 'USDC'], ['USDC', 'DAI', 'USDT'],
      ];
      for (const pair of stablePairs) {
        if (Math.random() > 0.6) continue;
        const symbol = pair.join('-');
        const baseApy = random(1, 15) * chain.apyMultiplier;
        const rewardApy = Math.random() > 0.4 ? random(2, 20) : 0;
        const tvl = random(500000, 200000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, symbol),
          chain: chain.name,
          protocol: protocol.name,
          symbol,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || [],
          underlyingTokens: pair,
          poolMeta: `${protocol.displayName} ${symbol} Pool`,
          il7d: random(0, 0.1),
          apyMean30d: baseApy + rewardApy + random(-2, 2),
          volumeUsd1d: tvl * random(0.05, 0.3),
          volumeUsd7d: tvl * random(0.2, 1.5),
          score: random(55, 90, 0),
          apyChange1h: random(-1, 1),
          apyChange24h: random(-5, 5),
          apyChange7d: random(-10, 10),
          stablecoin: true,
          exposure: 'multi',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }

      // Volatile pairs
      const volatilePairs = [
        ['ETH', 'USDC'], ['ETH', 'USDT'], ['WBTC', 'ETH'], ['ETH', 'DAI'],
        ['ARB', 'ETH'], ['OP', 'ETH'], ['MATIC', 'ETH'], ['LINK', 'ETH'],
        ['UNI', 'ETH'], ['AAVE', 'ETH'], ['CRV', 'ETH'], ['LDO', 'ETH'],
        ['GMX', 'ETH'], ['RDNT', 'ETH'], ['PENDLE', 'ETH'], ['stETH', 'ETH'],
        ['wstETH', 'ETH'], ['rETH', 'ETH'], ['cbETH', 'ETH'],
      ];
      for (const pair of volatilePairs) {
        if (Math.random() > 0.7) continue;
        const symbol = pair.join('-');
        const baseApy = random(5, 50) * chain.apyMultiplier;
        const rewardApy = Math.random() > 0.3 ? random(5, 40) : 0;
        const tvl = random(100000, 100000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, symbol),
          chain: chain.name,
          protocol: protocol.name,
          symbol,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || [],
          underlyingTokens: pair,
          poolMeta: `${protocol.displayName} ${symbol} Pool`,
          il7d: random(0.5, 5),
          apyMean30d: baseApy + rewardApy + random(-5, 5),
          volumeUsd1d: tvl * random(0.1, 0.8),
          volumeUsd7d: tvl * random(0.5, 4),
          score: random(45, 85, 0),
          apyChange1h: random(-2, 2),
          apyChange24h: random(-10, 10),
          apyChange7d: random(-20, 20),
          stablecoin: false,
          exposure: 'multi',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }
    }
  }

  // Generate liquid staking pools
  const stakingProtocols = protocols.filter(p => p.category === 'liquid-staking' || p.category === 'restaking');
  for (const protocol of stakingProtocols) {
    for (const chain of chains.slice(0, 5)) { // Mainly on major chains
      if (Math.random() > 0.4) continue;

      for (const token of tokens.lsts.slice(0, Math.floor(Math.random() * 5) + 2)) {
        if (Math.random() > 0.5) continue;
        const baseApy = random(3, 8);
        const rewardApy = Math.random() > 0.6 ? random(1, 10) : 0;
        const tvl = random(10000000, 5000000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, token),
          chain: chain.name,
          protocol: protocol.name,
          symbol: token,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || [],
          underlyingTokens: ['ETH'],
          poolMeta: `${protocol.displayName} ${token}`,
          il7d: 0,
          apyMean30d: baseApy + rewardApy + random(-0.3, 0.3),
          volumeUsd1d: tvl * random(0.005, 0.05),
          volumeUsd7d: tvl * random(0.02, 0.2),
          score: random(75, 98, 0),
          apyChange1h: random(-0.1, 0.1),
          apyChange24h: random(-0.5, 0.5),
          apyChange7d: random(-1, 1),
          stablecoin: false,
          exposure: 'single',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }
    }
  }

  // Generate yield aggregator pools
  const yieldProtocols = protocols.filter(p => p.category === 'yield');
  for (const protocol of yieldProtocols) {
    for (const chain of chains) {
      if (Math.random() > 0.6) continue;

      // Auto-compounding vaults
      const vaultAssets = [...tokens.stablecoins.slice(0, 5), ...tokens.majors.slice(0, 4), ...tokens.lsts.slice(0, 3)];
      for (const token of vaultAssets) {
        if (Math.random() > 0.75) continue;
        const baseApy = random(5, 25) * chain.apyMultiplier;
        const rewardApy = Math.random() > 0.5 ? random(2, 15) : 0;
        const tvl = random(100000, 50000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, `${token}-vault`),
          chain: chain.name,
          protocol: protocol.name,
          symbol: `${token} Vault`,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || [],
          underlyingTokens: [token],
          poolMeta: `${protocol.displayName} ${token} Auto-compound Vault`,
          il7d: tokens.stablecoins.includes(token) ? 0 : random(0, 2),
          apyMean30d: baseApy + rewardApy + random(-3, 3),
          volumeUsd1d: tvl * random(0.02, 0.15),
          volumeUsd7d: tvl * random(0.1, 0.8),
          score: random(50, 88, 0),
          apyChange1h: random(-1, 1),
          apyChange24h: random(-5, 5),
          apyChange7d: random(-12, 12),
          stablecoin: tokens.stablecoins.includes(token),
          exposure: 'single',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }
    }
  }

  // Generate derivatives pools (GMX, Gains, etc.)
  const derivativeProtocols = protocols.filter(p => p.category === 'derivatives');
  for (const protocol of derivativeProtocols) {
    for (const chain of chains.slice(0, 8)) {
      if (Math.random() > 0.5) continue;

      const derivativeProducts = ['GLP', 'GM-ETH', 'GM-BTC', 'GM-LINK', 'GM-ARB', 'gDAI', 'gUSDC'];
      for (const product of derivativeProducts) {
        if (Math.random() > 0.6) continue;
        const baseApy = random(10, 40) * chain.apyMultiplier;
        const rewardApy = random(5, 25);
        const tvl = random(5000000, 500000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, product),
          chain: chain.name,
          protocol: protocol.name,
          symbol: product,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || ['ETH'],
          underlyingTokens: product.includes('ETH') ? ['ETH'] : product.includes('BTC') ? ['BTC'] : ['USDC', 'ETH', 'BTC'],
          poolMeta: `${protocol.displayName} ${product}`,
          il7d: random(0.2, 3),
          apyMean30d: baseApy + rewardApy + random(-8, 8),
          volumeUsd1d: tvl * random(0.2, 1),
          volumeUsd7d: tvl * random(1, 5),
          score: random(55, 85, 0),
          apyChange1h: random(-2, 2),
          apyChange24h: random(-8, 8),
          apyChange7d: random(-15, 15),
          stablecoin: product.includes('USDC') || product.includes('DAI'),
          exposure: 'multi',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }
    }
  }

  // Generate bridge pools
  const bridgeProtocols = protocols.filter(p => p.category === 'bridge');
  for (const protocol of bridgeProtocols) {
    for (const chain of chains) {
      if (Math.random() > 0.5) continue;

      const bridgeAssets = ['USDC', 'USDT', 'ETH', 'DAI'];
      for (const token of bridgeAssets) {
        if (Math.random() > 0.6) continue;
        const baseApy = random(2, 12) * chain.apyMultiplier;
        const rewardApy = Math.random() > 0.4 ? random(1, 8) : 0;
        const tvl = random(1000000, 100000000) * chain.tvlMultiplier;

        pools.push({
          id: generatePoolId(protocol.name, chain.name, token),
          chain: chain.name,
          protocol: protocol.name,
          symbol: `${token} Bridge LP`,
          tvl,
          apy: baseApy + rewardApy,
          apyBase: baseApy,
          apyReward: rewardApy,
          rewardTokens: rewardTokens[protocol.name] || [],
          underlyingTokens: [token],
          poolMeta: `${protocol.displayName} ${token} Liquidity`,
          il7d: 0,
          apyMean30d: baseApy + rewardApy + random(-1, 1),
          volumeUsd1d: tvl * random(0.1, 0.5),
          volumeUsd7d: tvl * random(0.5, 2),
          score: random(60, 90, 0),
          apyChange1h: random(-0.5, 0.5),
          apyChange24h: random(-2, 2),
          apyChange7d: random(-5, 5),
          stablecoin: token !== 'ETH',
          exposure: 'single',
          createdAt: new Date(now.getTime() - random(30, 365) * 86400000).toISOString(),
          updatedAt: now.toISOString(),
        });
      }
    }
  }

  return pools;
}

// Generate all mock pools
export const mockPools: Pool[] = generatePools();

// Generate opportunities based on pools
function generateOpportunities(): Opportunity[] {
  const opportunities: Opportunity[] = [];
  const now = new Date();

  // Find yield gaps (same asset, different protocols)
  const poolsByAsset: Record<string, Pool[]> = {};
  for (const pool of mockPools) {
    const key = `${pool.chain}-${pool.symbol}`;
    if (!poolsByAsset[key]) poolsByAsset[key] = [];
    poolsByAsset[key].push(pool);
  }

  for (const [key, pools] of Object.entries(poolsByAsset)) {
    if (pools.length < 2) continue;
    pools.sort((a, b) => b.apy - a.apy);
    const best = pools[0];
    const worst = pools[pools.length - 1];
    const diff = best.apy - worst.apy;

    if (diff > 1) {
      opportunities.push({
        id: `yield-gap-${key}`,
        type: 'yield-gap',
        title: `${best.symbol} Yield Gap: ${best.protocol} vs ${worst.protocol}`,
        description: `${best.protocol} offers ${diff.toFixed(2)}% higher APY on ${best.symbol} compared to ${worst.protocol}. Consider migrating for better returns.`,
        sourcePoolId: worst.id,
        targetPoolId: best.id,
        asset: best.symbol,
        chain: best.chain,
        apyDifference: diff,
        apyGrowth: 0,
        currentApy: best.apy,
        potentialProfit: diff * 100,
        tvl: best.tvl,
        riskLevel: diff > 5 ? 'medium' : 'low',
        score: Math.min(95, 60 + diff * 3),
        isActive: true,
        detectedAt: new Date(now.getTime() - random(1, 24) * 3600000).toISOString(),
        lastSeenAt: now.toISOString(),
        expiresAt: new Date(now.getTime() + random(12, 72) * 3600000).toISOString(),
        createdAt: new Date(now.getTime() - random(1, 24) * 3600000).toISOString(),
        updatedAt: now.toISOString(),
      });
    }
  }

  // Find trending pools (high APY growth)
  const trendingPools = mockPools
    .filter(p => p.apyChange24h > 5)
    .sort((a, b) => b.apyChange24h - a.apyChange24h)
    .slice(0, 30);

  for (const pool of trendingPools) {
    opportunities.push({
      id: `trending-${pool.id}`,
      type: 'trending',
      title: `${pool.symbol} APY Surging +${pool.apyChange24h.toFixed(1)}% on ${pool.protocol}`,
      description: `${pool.protocol} ${pool.symbol} pool showing strong APY growth. Current APY: ${pool.apy.toFixed(2)}%.`,
      poolId: pool.id,
      asset: pool.symbol,
      chain: pool.chain,
      apyDifference: 0,
      apyGrowth: pool.apyChange24h,
      currentApy: pool.apy,
      potentialProfit: pool.apy * 100,
      tvl: pool.tvl,
      riskLevel: pool.apyChange24h > 20 ? 'high' : pool.apyChange24h > 10 ? 'medium' : 'low',
      score: Math.min(90, 50 + pool.apyChange24h * 2),
      isActive: true,
      detectedAt: new Date(now.getTime() - random(0.5, 12) * 3600000).toISOString(),
      lastSeenAt: now.toISOString(),
      expiresAt: new Date(now.getTime() + random(6, 48) * 3600000).toISOString(),
      createdAt: new Date(now.getTime() - random(0.5, 12) * 3600000).toISOString(),
      updatedAt: now.toISOString(),
    });
  }

  // Find high-score opportunities
  const highScorePools = mockPools
    .filter(p => p.score >= 85)
    .sort((a, b) => b.score - a.score)
    .slice(0, 20);

  for (const pool of highScorePools) {
    opportunities.push({
      id: `high-score-${pool.id}`,
      type: 'high-score',
      title: `${pool.symbol} - Premium Risk-Adjusted Yield (Score: ${pool.score})`,
      description: `${pool.protocol} ${pool.symbol} offers excellent risk-adjusted returns with ${pool.apy.toFixed(2)}% APY and $${(pool.tvl / 1e6).toFixed(0)}M TVL.`,
      poolId: pool.id,
      asset: pool.symbol,
      chain: pool.chain,
      apyDifference: 0,
      apyGrowth: pool.apyChange24h,
      currentApy: pool.apy,
      potentialProfit: pool.apy * 100,
      tvl: pool.tvl,
      riskLevel: 'low',
      score: pool.score,
      isActive: true,
      detectedAt: new Date(now.getTime() - random(2, 48) * 3600000).toISOString(),
      lastSeenAt: now.toISOString(),
      expiresAt: new Date(now.getTime() + random(24, 168) * 3600000).toISOString(),
      createdAt: new Date(now.getTime() - random(2, 48) * 3600000).toISOString(),
      updatedAt: now.toISOString(),
    });
  }

  return opportunities.slice(0, 100); // Limit to 100 opportunities
}

export const mockOpportunities: Opportunity[] = generateOpportunities();

// Calculate stats from pools
function calculateStats(): PlatformStats {
  const totalTvl = mockPools.reduce((sum, p) => sum + p.tvl, 0);
  const apys = mockPools.map(p => p.apy).sort((a, b) => a - b);
  const avgApy = apys.reduce((sum, a) => sum + a, 0) / apys.length;
  const medianApy = apys[Math.floor(apys.length / 2)];
  const maxApy = Math.max(...apys);

  const tvlByChain: Record<string, number> = {};
  const poolsByChain: Record<string, number> = {};
  const uniqueProtocols = new Set<string>();

  for (const pool of mockPools) {
    tvlByChain[pool.chain] = (tvlByChain[pool.chain] || 0) + pool.tvl;
    poolsByChain[pool.chain] = (poolsByChain[pool.chain] || 0) + 1;
    uniqueProtocols.add(pool.protocol);
  }

  const apyDistribution = {
    range0to1: mockPools.filter(p => p.apy < 1).length,
    range1to5: mockPools.filter(p => p.apy >= 1 && p.apy < 5).length,
    range5to10: mockPools.filter(p => p.apy >= 5 && p.apy < 10).length,
    range10to25: mockPools.filter(p => p.apy >= 10 && p.apy < 25).length,
    range25to50: mockPools.filter(p => p.apy >= 25 && p.apy < 50).length,
    range50to100: mockPools.filter(p => p.apy >= 50 && p.apy < 100).length,
    range100plus: mockPools.filter(p => p.apy >= 100).length,
  };

  return {
    totalPools: mockPools.length,
    totalTvl,
    averageApy: avgApy,
    medianApy,
    maxApy,
    totalChains: Object.keys(tvlByChain).length,
    totalProtocols: uniqueProtocols.size,
    activeOpportunities: mockOpportunities.filter(o => o.isActive).length,
    lastUpdated: new Date().toISOString(),
    tvlByChain,
    poolsByChain,
    apyDistribution,
  };
}

export const mockStats: PlatformStats = calculateStats();

// Generate chains data from pools
function generateChains(): Chain[] {
  const chainMap: Record<string, { pools: Pool[]; tvl: number }> = {};

  for (const pool of mockPools) {
    if (!chainMap[pool.chain]) {
      chainMap[pool.chain] = { pools: [], tvl: 0 };
    }
    chainMap[pool.chain].pools.push(pool);
    chainMap[pool.chain].tvl += pool.tvl;
  }

  return Object.entries(chainMap)
    .map(([name, data]) => {
      const apys = data.pools.map(p => p.apy);
      const protocolCounts: Record<string, number> = {};
      for (const p of data.pools) {
        protocolCounts[p.protocol] = (protocolCounts[p.protocol] || 0) + 1;
      }
      const topProtocols = Object.entries(protocolCounts)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 5)
        .map(([p]) => p);

      const chainInfo = chains.find(c => c.name === name);

      return {
        name,
        displayName: chainInfo?.displayName || name,
        poolCount: data.pools.length,
        totalTvl: data.tvl,
        averageApy: apys.reduce((s, a) => s + a, 0) / apys.length,
        maxApy: Math.max(...apys),
        topProtocols,
      };
    })
    .sort((a, b) => b.totalTvl - a.totalTvl);
}

export const mockChains: Chain[] = generateChains();

// Generate trending pools
function generateTrendingPools(): TrendingPool[] {
  return mockPools
    .filter(p => p.apyChange24h > 3)
    .sort((a, b) => b.apyChange24h - a.apyChange24h)
    .slice(0, 20)
    .map(pool => ({
      pool,
      apyGrowth1h: pool.apyChange1h,
      apyGrowth24h: pool.apyChange24h,
      apyGrowth7d: pool.apyChange7d,
      trendScore: Math.min(100, 50 + pool.apyChange24h * 2.5),
    }));
}

export const mockTrendingPools: TrendingPool[] = generateTrendingPools();

// Generate mock historical data
export function generateMockHistory(poolId: string, period: '1h' | '24h' | '7d' | '30d') {
  const now = Date.now();
  const points: { poolId: string; timestamp: string; apy: number; tvl: number; apyBase: number; apyReward: number }[] = [];

  let intervalMs: number;
  let numPoints: number;

  switch (period) {
    case '1h':
      intervalMs = 60000;
      numPoints = 60;
      break;
    case '24h':
      intervalMs = 3600000;
      numPoints = 24;
      break;
    case '7d':
      intervalMs = 86400000 / 4;
      numPoints = 28;
      break;
    case '30d':
      intervalMs = 86400000;
      numPoints = 30;
      break;
  }

  const pool = mockPools.find(p => p.id === poolId) || mockPools[0];
  const baseApy = pool.apy;
  const baseTvl = pool.tvl;

  for (let i = numPoints - 1; i >= 0; i--) {
    const timestamp = new Date(now - i * intervalMs).toISOString();
    const variance = (Math.random() - 0.5) * 4;
    const tvlVariance = (Math.random() - 0.5) * 0.15;

    points.push({
      poolId,
      timestamp,
      apy: Math.max(0, baseApy + variance),
      tvl: baseTvl * (1 + tvlVariance),
      apyBase: Math.max(0, pool.apyBase + variance * 0.6),
      apyReward: Math.max(0, pool.apyReward + variance * 0.4),
    });
  }

  return { poolId, period, dataPoints: points };
}

// Generated pools: mockPools.length across mockChains.length chains
