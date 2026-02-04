# Recurring Swaps

## Summary
Automate token swaps on a schedule (DCA - Dollar Cost Averaging).

## Capabilities
- **Schedule recurring swaps**: Convert one token to another on a daily, weekly, or monthly basis
- **Cross-chain swaps**: Swap tokens across supported chains via THORChain and MayaChain DEXs
- **Same-chain swaps**: Swap tokens on the same chain via DEX aggregators

## Supported Chains
Ethereum, Arbitrum, Avalanche, BNB Chain, Base, Blast, Optimism, Polygon, Bitcoin, Litecoin, Dogecoin, Bitcoin Cash, Dash, Solana, XRP, Zcash, Cosmos (Gaia), MayaChain, THORChain, Tron

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| from.chain | Yes | Source blockchain |
| from.token | No | Token address or empty for native asset |
| from.address | Yes | User's wallet address |
| to.chain | Yes | Destination blockchain |
| to.token | No | Token address or empty for native asset |
| to.address | Yes | Destination wallet address |
| fromAmount | Yes | Amount in smallest unit (wei, satoshi, etc.) |
| frequency | Yes | "one-time", "minutely", "hourly", "daily", "weekly", "bi-weekly", or "monthly" |
| startDate | No | ISO 8601 datetime when to start (default: now) |
| endDate | No | ISO 8601 datetime when to stop |

## Example User Requests
- "Set up a weekly DCA into ETH"
- "Swap $100 USDC to ETH every week"
- "I want to dollar cost average into Bitcoin"
- "Convert my USDC to ETH automatically"
- "Buy 0.01 BTC with USDC every day"
- "Swap AVAX to ETH on Arbitrum monthly"

## Limitations
- Minimum swap amounts apply per chain
- Cross-chain swaps execute via THORChain or MayaChain DEX
- Not all token pairs are supported - destination asset must have liquidity in the DEX
- Swaps require sufficient balance for the swap amount plus gas fees
