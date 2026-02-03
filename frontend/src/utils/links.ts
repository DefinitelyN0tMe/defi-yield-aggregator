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
  'fantom': 'https://ftmscan.com',
  'gnosis': 'https://gnosisscan.io',
  'zkSync': 'https://explorer.zksync.io',
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
};

export const defiLlamaLinks = {
  pool: (poolId: string) => `https://defillama.com/yields/pool/${poolId}`,
  protocol: (protocol: string) => `https://defillama.com/protocol/${protocol}`,
  chain: (chain: string) => `https://defillama.com/chain/${chain}`,
};

export function getProtocolUrl(protocol: string, chain?: string): string {
  const baseUrl = protocolLinks[protocol];
  if (!baseUrl) return `https://defillama.com/protocol/${protocol}`;

  // Some protocols have chain-specific URLs
  if (chain && baseUrl.includes('app.')) {
    // Could add chain-specific routing here if needed
  }

  return baseUrl;
}

export function getExplorerUrl(chain: string): string {
  return chainExplorers[chain] || `https://defillama.com/chain/${chain}`;
}

export function getTokenUrl(chain: string, tokenAddress?: string): string {
  const explorer = chainExplorers[chain];
  if (!explorer || !tokenAddress) return '';
  return `${explorer}/token/${tokenAddress}`;
}
