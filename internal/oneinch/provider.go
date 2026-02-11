package oneinch

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	evm_swap "github.com/vultisig/app-recurring/internal/evm"
	"github.com/vultisig/recipes/sdk/evm"
)

type Provider struct {
	client *Client
	rpc    *ethclient.Client
	sdk    *evm.SDK
}

func NewProvider(client *Client, rpc *ethclient.Client, sdk *evm.SDK) *Provider {
	return &Provider{
		client: client,
		rpc:    rpc,
		sdk:    sdk,
	}
}

func (p *Provider) Name() string {
	return "1inch"
}

func (p *Provider) validateSwap(from evm_swap.From, to evm_swap.To) error {
	if !from.Chain.IsEvm() {
		return fmt.Errorf("from chain %s is not EVM", from.Chain.String())
	}

	if !to.Chain.IsEvm() {
		return fmt.Errorf("to chain %s is not EVM", to.Chain.String())
	}

	if from.Chain != to.Chain {
		return fmt.Errorf("1inch only supports same-chain swaps: from=%s, to=%s", from.Chain.String(), to.Chain.String())
	}

	return nil
}

type resolvedSwap struct {
	resp     *SwapResponse
	toAmount *big.Int
}

func (p *Provider) resolveSwap(
	ctx context.Context,
	from evm_swap.From,
	to evm_swap.To,
) (*resolvedSwap, error) {
	if err := p.validateSwap(from, to); err != nil {
		return nil, fmt.Errorf("invalid swap: %w", err)
	}

	srcToken := from.AssetID.Hex()
	if bytes.Equal(from.AssetID.Bytes(), evm.ZeroAddress.Bytes()) {
		srcToken = NativeTokenAddress
	}

	dstToken := to.AssetID
	if to.AssetID == "" {
		dstToken = NativeTokenAddress
	}

	swapResp, err := p.client.GetSwap(ctx, swapRequest{
		Chain:        from.Chain,
		Src:          srcToken,
		Dst:          dstToken,
		Amount:       from.Amount.String(),
		From:         from.Address.Hex(),
		Receiver:     to.Address,
		SlippagePerc: DefaultSlippagePercent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get swap from 1inch: %w", err)
	}

	toAmount, ok := new(big.Int).SetString(swapResp.DstAmount, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse toAmount: %s", swapResp.DstAmount)
	}

	return &resolvedSwap{
		resp:     swapResp,
		toAmount: toAmount,
	}, nil
}

func (p *Provider) Quote(
	ctx context.Context,
	from evm_swap.From,
	to evm_swap.To,
) (*evm_swap.QuoteResult, error) {
	rs, err := p.resolveSwap(ctx, from, to)
	if err != nil {
		return nil, err
	}

	spenderAddr, err := p.client.GetSpender(ctx, from.Chain)
	if err != nil {
		return nil, fmt.Errorf("failed to get spender: %w", err)
	}

	return &evm_swap.QuoteResult{
		AmountOut: rs.toAmount,
		Spender:   ecommon.HexToAddress(spenderAddr),
	}, nil
}

func (p *Provider) MakeTx(
	ctx context.Context,
	from evm_swap.From,
	to evm_swap.To,
) (*big.Int, []byte, error) {
	rs, err := p.resolveSwap(ctx, from, to)
	if err != nil {
		return nil, nil, err
	}

	txValue, ok := new(big.Int).SetString(rs.resp.Tx.Value, 10)
	if !ok {
		return nil, nil, fmt.Errorf("failed to parse tx value: %s", rs.resp.Tx.Value)
	}

	txData, err := hexutil.Decode(rs.resp.Tx.Data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode tx data: %w", err)
	}

	routerAddr := ecommon.HexToAddress(rs.resp.Tx.To)

	unsignedTx, err := p.sdk.MakeTx(
		ctx,
		from.Address,
		routerAddr,
		txValue,
		txData,
		0, // nonceOffset: swaps don't need offset
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build unsigned tx: %w", err)
	}

	return rs.toAmount, unsignedTx, nil
}
