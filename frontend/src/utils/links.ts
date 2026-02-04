// External links to DeFi protocols and explorers

export const protocolLinks: Record<string, string> = {
  'aave-v3': 'https://app.aave.com/',
  'aave-v2': 'https://app.aave.com/',
  'compound-v3': 'https://app.compound.finance/',
  'compound-v2': 'https://app.compound.finance/',
  'lido': 'https://stake.lido.fi/',
  'rocket-pool': 'https://stake.rocketpool.net/',
  'frax-ether': 'https://app.frax.finance/',
  'curve': 'https://curve.fi/',
  'convex': 'https://www.convexfinance.com/',
  'uniswap-v3': 'https://app.uniswap.org/',
  'uniswap-v2': 'https://app.uniswap.org/',
  'sushiswap': 'https://www.sushi.com/',
  'balancer-v2': 'https://app.balancer.fi/',
  'gmx': 'https://app.gmx.io/',
  'gains-network': 'https://gains.trade/',
  'pendle': 'https://app.pendle.finance/',
  'yearn': 'https://yearn.fi/',
  'beefy': 'https://app.beefy.com/',
  'velodrome': 'https://velodrome.finance/',
  'aerodrome': 'https://aerodrome.finance/',
  'camelot': 'https://app.camelot.exchange/',
  'trader-joe': 'https://traderjoexyz.com/',
  'pancakeswap': 'https://pancakeswap.finance/',
  'quickswap': 'https://quickswap.exchange/',
  'radiant': 'https://app.radiant.capital/',
  'morpho': 'https://app.morpho.org/',
  'spark': 'https://app.spark.fi/',
  'venus': 'https://app.venus.io/',
  'benqi': 'https://app.benqi.fi/',
  'stargate': 'https://stargate.finance/',
  'hop-protocol': 'https://app.hop.exchange/',
  'across': 'https://across.to/',
  'eigenlayer': 'https://app.eigenlayer.xyz/',
  'ether-fi': 'https://app.ether.fi/',
  'renzo': 'https://app.renzoprotocol.com/',
  'kelp-dao': 'https://kelpdao.xyz/',
  'sommelier': 'https://app.sommelier.finance/',
  'instadapp': 'https://instadapp.io/',
  'maple': 'https://app.maple.finance/',
  'goldfinch': 'https://app.goldfinch.finance/',
};

export const chainExplorers: Record<string, string> = {
  'ethereum': 'https://etherscan.io',
  'arbitrum': 'https://arbiscan.io',
  'optimism': 'https://optimistic.etherscan.io',
  'polygon': 'https://polygonscan.com',
  'base': 'https://basescan.org',
  'bsc': 'https://bscscan.com',
  'avalanche': 'https://snowtrace.io',
  'fantom': 'https://www.oklink.com/fantom',
  'gnosis': 'https://gnosisscan.io',
  'zksync era': 'https://explorer.zksync.io',
  'linea': 'https://lineascan.build',
  'scroll': 'https://scrollscan.com',
  'mantle': 'https://explorer.mantle.xyz',
  'manta': 'https://pacific-explorer.manta.network',
  'blast': 'https://blastscan.io',
  'mode': 'https://explorer.mode.network',
  'metis': 'https://andromeda-explorer.metis.io',
  'celo': 'https://celoscan.io',
  'moonbeam': 'https://moonscan.io',
  'kava': 'https://kavascan.com',
  'solana': 'https://solscan.io',
  'sui': 'https://suiscan.xyz',
  'aptos': 'https://explorer.aptoslabs.com',
  'tron': 'https://tronscan.org',
  'cronos': 'https://cronoscan.com',
  'aurora': 'https://aurorascan.dev',
  'moonriver': 'https://moonriver.moonscan.io',
  'harmony': 'https://explorer.harmony.one',
  'near': 'https://nearblocks.io',
  'sei': 'https://seitrace.com',
  'ton': 'https://tonscan.org',
  'starknet': 'https://starkscan.co',
  'berachain': 'https://berascan.com',
  'sonic': 'https://sonicscan.org',
  'hyperliquid': 'https://hyperliquid.xyz',
  'monad': 'https://monad.xyz',
};

// CoinGecko ID mapping for common DeFi tokens
// CoinGecko uses slug IDs, not symbols
const coinGeckoIds: Record<string, string> = {
  // Stablecoins
  'USDC': 'usd-coin',
  'USDT': 'tether',
  'DAI': 'dai',
  'FRAX': 'frax',
  'LUSD': 'liquity-usd',
  'crvUSD': 'crvusd',
  'GHO': 'gho',
  'PYUSD': 'paypal-usd',
  'TUSD': 'true-usd',
  'BUSD': 'binance-usd',
  'USDP': 'paxos-standard',
  'sUSD': 'susd',
  'MIM': 'magic-internet-money',
  'DOLA': 'dola-usd',
  'alUSD': 'alchemix-usd',

  // ETH and LSTs
  'ETH': 'ethereum',
  'WETH': 'weth',
  'stETH': 'lido-staked-ether',
  'wstETH': 'wrapped-steth',
  'rETH': 'rocket-pool-eth',
  'cbETH': 'coinbase-wrapped-staked-eth',
  'frxETH': 'frax-ether',
  'sfrxETH': 'staked-frax-ether',
  'ankrETH': 'ankr-staked-eth',
  'swETH': 'sweth',
  'ETHx': 'stader-ethx',
  'mETH': 'mantle-staked-ether',
  'eETH': 'ether-fi-staked-eth',
  'weETH': 'wrapped-eeth',
  'ezETH': 'renzo-restaked-eth',
  'pufETH': 'puffer-finance',
  'rsETH': 'kelp-dao-restaked-eth',

  // BTC
  'WBTC': 'wrapped-bitcoin',
  'BTC': 'bitcoin',
  'tBTC': 'tbtc',

  // Major DeFi tokens
  'LINK': 'chainlink',
  'UNI': 'uniswap',
  'AAVE': 'aave',
  'CRV': 'curve-dao-token',
  'CVX': 'convex-finance',
  'LDO': 'lido-dao',
  'RPL': 'rocket-pool',
  'GMX': 'gmx',
  'ARB': 'arbitrum',
  'OP': 'optimism',
  'MATIC': 'matic-network',
  'SNX': 'havven',
  'COMP': 'compound-governance-token',
  'MKR': 'maker',
  'BAL': 'balancer',
  'SUSHI': 'sushi',
  'YFI': 'yearn-finance',
  'FXS': 'frax-share',
  'SPELL': 'spell-token',
  'ALCX': 'alchemix',
  'PENDLE': 'pendle',
  'VELO': 'velodrome-finance',
  'AERO': 'aerodrome-finance',
  'CAKE': 'pancakeswap-token',
  'JOE': 'joe',
  'QUICK': 'quickswap',
  'RDNT': 'radiant-capital',
  'STG': 'stargate-finance',
  'GRAIL': 'camelot-token',
  'BIFI': 'beefy-finance',

  // Chain native tokens
  'BNB': 'binancecoin',
  'AVAX': 'avalanche-2',
  'FTM': 'fantom',
  'CELO': 'celo',
  'GLMR': 'moonbeam',
  'KAVA': 'kava',
  'METIS': 'metis-token',
  'MANTA': 'manta-network',
  'MODE': 'mode',
};

// DefiLlama protocol ID mapping (some protocols have different IDs)
const defiLlamaProtocolIds: Record<string, string> = {
  'aave-v3': 'aave-v3',
  'aave-v2': 'aave',
  'compound-v3': 'compound-v3',
  'compound-v2': 'compound',
  'lido': 'lido',
  'rocket-pool': 'rocket-pool',
  'curve': 'curve-finance',
  'convex': 'convex-finance',
  'uniswap-v3': 'uniswap-v3',
  'uniswap-v2': 'uniswap',
  'sushiswap': 'sushi',
  'balancer-v2': 'balancer-v2',
  'gmx': 'gmx',
  'pendle': 'pendle',
  'yearn': 'yearn-finance',
  'beefy': 'beefy',
  'velodrome': 'velodrome',
  'aerodrome': 'aerodrome',
  'pancakeswap': 'pancakeswap',
  'radiant': 'radiant-v2',
  'morpho': 'morpho',
  'spark': 'spark',
  'stargate': 'stargate',
  'eigenlayer': 'eigenlayer',
  'ether-fi': 'ether.fi',
  'renzo': 'renzo',
};

export const defiLlamaLinks = {
  pool: (poolId: string) => `https://defillama.com/yields/pool/${poolId}`,
  protocol: (protocol: string) => {
    const id = defiLlamaProtocolIds[protocol] || protocol;
    return `https://defillama.com/protocol/${id}`;
  },
  chain: (chain: string) => `https://defillama.com/chain/${chain}`,
  yields: () => 'https://defillama.com/yields',
};

export function getProtocolUrl(protocol: string, chain?: string): string {
  const baseUrl = protocolLinks[protocol];
  if (!baseUrl) return defiLlamaLinks.protocol(protocol);

  // Some protocols have chain-specific URLs
  if (chain && baseUrl.includes('app.')) {
    // Could add chain-specific routing here if needed
  }

  return baseUrl;
}

export function getExplorerUrl(chain: string): string {
  const lowerChain = chain.toLowerCase();
  return chainExplorers[lowerChain] || `https://defillama.com/chain/${chain}`;
}

export function getExplorerTokenSearchUrl(chain: string, tokenSymbol: string): string {
  const lowerChain = chain.toLowerCase();
  const explorer = chainExplorers[lowerChain];
  if (!explorer) return `https://defillama.com/chain/${chain}`;

  // Most explorers support token search
  // Etherscan-like explorers use /tokens?q=
  if (explorer.includes('etherscan') || explorer.includes('scan')) {
    return `${explorer}/tokens?q=${encodeURIComponent(tokenSymbol)}`;
  }

  // zkSync uses different search
  if (explorer.includes('zksync')) {
    return `${explorer}/tokens`;
  }

  // Fallback to explorer homepage
  return explorer;
}

export function getTokenUrl(chain: string, tokenAddress?: string): string {
  const lowerChain = chain.toLowerCase();
  const explorer = chainExplorers[lowerChain];
  if (!explorer || !tokenAddress) return '';
  return `${explorer}/token/${tokenAddress}`;
}

export function getCoinGeckoUrl(tokenSymbol: string): string {
  // Clean up the symbol - remove wrapping prefixes and get base token
  const cleanSymbol = tokenSymbol.toUpperCase().trim();

  // Try exact match first
  if (coinGeckoIds[cleanSymbol]) {
    return `https://www.coingecko.com/en/coins/${coinGeckoIds[cleanSymbol]}`;
  }

  // For LP tokens like "WETH-USDC", try the first token
  if (cleanSymbol.includes('-')) {
    const firstToken = cleanSymbol.split('-')[0];
    if (coinGeckoIds[firstToken]) {
      return `https://www.coingecko.com/en/coins/${coinGeckoIds[firstToken]}`;
    }
  }

  // For LP tokens like "WETH/USDC", try the first token
  if (cleanSymbol.includes('/')) {
    const firstToken = cleanSymbol.split('/')[0];
    if (coinGeckoIds[firstToken]) {
      return `https://www.coingecko.com/en/coins/${coinGeckoIds[firstToken]}`;
    }
  }

  // Fallback to search - CoinGecko's search is reliable
  return `https://www.coingecko.com/en/search?query=${encodeURIComponent(tokenSymbol)}`;
}

export function getDexScreenerUrl(tokenSymbol: string, chain?: string): string {
  // DexScreener search works well
  const cleanSymbol = tokenSymbol.split('-')[0].split('/')[0];
  if (chain) {
    return `https://dexscreener.com/${chain.toLowerCase()}?q=${encodeURIComponent(cleanSymbol)}`;
  }
  return `https://dexscreener.com/search?q=${encodeURIComponent(cleanSymbol)}`;
}
