# Recurring Sends

## Summary
Automate token transfers on a schedule (recurring payments, payroll, subscriptions).

## Capabilities
- **Schedule recurring transfers**: Send tokens to one or more recipients on a regular schedule
- **Multi-recipient support**: Configure up to 10 recipients per policy
- **Native and token transfers**: Send both native assets (ETH, BTC, SOL) and tokens (ERC20, SPL)
- **Memo support**: Include memos for chains that support them (XRP, Cosmos)

## Supported Chains
Ethereum, Arbitrum, Avalanche, BNB Chain, Base, Blast, Optimism, Polygon, Bitcoin, Litecoin, Dogecoin, Bitcoin Cash, Dash, Solana, XRP, Zcash, Cosmos (Gaia), MayaChain, THORChain, Tron

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| asset.chain | Yes | Blockchain to send on |
| asset.token | No | Token address or empty for native asset |
| asset.address | Yes | Sender's wallet address |
| recipients | Yes | Array of recipient objects (1-10 recipients) |
| recipients[].toAddress | Yes | Recipient wallet address |
| recipients[].amount | Yes | Amount to send in smallest unit |
| recipients[].alias | No | Optional friendly name for recipient |
| memo | No | Optional memo/tag for the transaction |
| frequency | Yes | "one-time", "minutely", "hourly", "daily", "weekly", "bi-weekly", or "monthly" |
| startDate | No | ISO 8601 datetime when to start (default: now) |
| endDate | No | ISO 8601 datetime when to stop |

## Example User Requests
- "Send 100 USDC to my friend every week"
- "Set up monthly payments of 0.5 ETH"
- "Pay my team 500 USDC each on the 1st of every month"
- "Transfer 0.01 BTC to savings wallet daily"
- "Send SOL to three addresses weekly"
- "Set up recurring XRP payment with memo tag"

## Limitations
- Maximum 10 recipients per policy
- Requires sufficient balance for all transfers plus gas fees
- Memo parameter is optional and only used by chains that support it
- Transfers are same-chain only (use Recurring Swaps for cross-chain)
