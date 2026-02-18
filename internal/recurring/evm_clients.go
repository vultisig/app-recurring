package recurring

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vultisig/recipes/sdk/evm"
	"github.com/vultisig/recipes/sdk/swap"
)

type evmClientPool struct {
	mu      sync.RWMutex
	clients map[string]*evm.SDK
}

var defaultEVMPool = &evmClientPool{
	clients: make(map[string]*evm.SDK),
}

func getEVMSDK(chain string) (*evm.SDK, error) {
	return defaultEVMPool.get(chain)
}

func (p *evmClientPool) get(chain string) (*evm.SDK, error) {
	p.mu.RLock()
	sdk, ok := p.clients[chain]
	p.mu.RUnlock()
	if ok {
		return sdk, nil
	}

	cfg, err := swap.GetEVMChainConfig(chain)
	if err != nil {
		return nil, fmt.Errorf("swap.GetEVMChainConfig(%s): %w", chain, err)
	}

	if len(cfg.RPCURLs) == 0 {
		return nil, fmt.Errorf("no RPC URLs for chain %s", chain)
	}

	rpcClient, err := ethclient.Dial(cfg.RPCURLs[0])
	if err != nil {
		return nil, fmt.Errorf("ethclient.Dial(%s): %w", cfg.RPCURLs[0], err)
	}

	sdk = evm.NewSDK(cfg.ChainID, rpcClient, rpcClient.Client())

	p.mu.Lock()
	p.clients[chain] = sdk
	p.mu.Unlock()

	return sdk, nil
}
