package recurring

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/vultisig/recipes/sdk/swap"
)

type BuildTxRequest struct {
	FromChain    string `json:"from_chain"`
	FromSymbol   string `json:"from_symbol"`
	FromAddress  string `json:"from_address,omitempty"`
	FromDecimals *int   `json:"from_decimals,omitempty"`
	ToChain      string `json:"to_chain"`
	ToSymbol     string `json:"to_symbol"`
	ToAddress    string `json:"to_address,omitempty"`
	ToDecimals   *int   `json:"to_decimals,omitempty"`
	Amount       string `json:"amount"`
	Sender       string `json:"sender"`
	Destination  string `json:"destination"`
}

type BuildTxResponse struct {
	Provider       string      `json:"provider"`
	ExpectedOutput string      `json:"expected_output"`
	MinimumOutput  string      `json:"minimum_output"`
	NeedsApproval  bool        `json:"needs_approval"`
	ApprovalTx     *TxDataJSON `json:"approval_tx,omitempty"`
	SwapTx         *TxDataJSON `json:"swap_tx"`
	Memo           string      `json:"memo,omitempty"`
}

type TxDataJSON struct {
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	Memo     string `json:"memo,omitempty"`
	Nonce    uint64 `json:"nonce"`
	GasLimit uint64 `json:"gas_limit"`
	ChainID  string `json:"chain_id,omitempty"`
}

func (s *SwapSpec) HandleBuildTx(ctx context.Context, body []byte) (any, error) {
	var req BuildTxRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	if req.FromChain == "" || req.FromSymbol == "" || req.ToChain == "" || req.ToSymbol == "" {
		return nil, fmt.Errorf("from_chain, from_symbol, to_chain, to_symbol are required")
	}
	if req.Amount == "" || req.Sender == "" || req.Destination == "" {
		return nil, fmt.Errorf("amount, sender, destination are required")
	}

	if req.FromDecimals == nil {
		return nil, fmt.Errorf("unknown token: %s on %s (from_decimals required)", req.FromSymbol, req.FromChain)
	}
	if req.ToDecimals == nil {
		return nil, fmt.Errorf("unknown token: %s on %s (to_decimals required)", req.ToSymbol, req.ToChain)
	}

	fromDecimals := *req.FromDecimals
	toDecimals := *req.ToDecimals

	if fromDecimals < 0 || toDecimals < 0 {
		return nil, fmt.Errorf("decimals must be non-negative")
	}

	amount, err := parseHumanAmount(req.Amount, fromDecimals)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	params := swap.SwapParams{
		FromChain:    req.FromChain,
		FromSymbol:   req.FromSymbol,
		FromAddress:  req.FromAddress,
		FromDecimals: fromDecimals,
		ToChain:      req.ToChain,
		ToSymbol:     req.ToSymbol,
		ToAddress:    req.ToAddress,
		ToDecimals:   toDecimals,
		Amount:       amount,
		Sender:       req.Sender,
		Destination:  req.Destination,
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	bundle, err := s.swapSvc.GetSwapTxBundle(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("swap provider error: %w", err)
	}

	expectedOutput := "0"
	if bundle.ExpectedOutput != nil {
		expectedOutput = bundle.ExpectedOutput.String()
	}
	minimumOutput := "0"
	if bundle.MinimumOutput != nil {
		minimumOutput = bundle.MinimumOutput.String()
	}

	resp := BuildTxResponse{
		Provider:       bundle.Provider,
		ExpectedOutput: expectedOutput,
		MinimumOutput:  minimumOutput,
		NeedsApproval:  bundle.NeedsApproval,
		Memo:           bundle.Memo,
	}

	if bundle.ApprovalTx != nil {
		resp.ApprovalTx = txDataToJSON(bundle.ApprovalTx)
	}
	if bundle.SwapTx != nil {
		resp.SwapTx = txDataToJSON(bundle.SwapTx)
	}

	return resp, nil
}

func txDataToJSON(td *swap.TxData) *TxDataJSON {
	result := &TxDataJSON{
		To:       td.To,
		Value:    "0",
		Data:     hex.EncodeToString(td.Data),
		Memo:     td.Memo,
		Nonce:    td.Nonce,
		GasLimit: td.GasLimit,
	}
	if td.Value != nil {
		result.Value = td.Value.String()
	}
	if td.ChainID != nil {
		result.ChainID = td.ChainID.String()
	}
	return result
}

func parseHumanAmount(s string, decimals int) (*big.Int, error) {
	if decimals < 0 {
		return nil, fmt.Errorf("decimals must be non-negative")
	}

	parts := strings.SplitN(s, ".", 2)
	wholePart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}

	if wholePart == "" && fracPart == "" {
		return nil, fmt.Errorf("cannot parse %q as number", s)
	}

	if len(fracPart) > decimals {
		fracPart = fracPart[:decimals]
	}
	for len(fracPart) < decimals {
		fracPart += "0"
	}

	combined := wholePart + fracPart
	if combined == "" {
		return nil, fmt.Errorf("cannot parse %q as number", s)
	}

	result := new(big.Int)
	_, ok := result.SetString(combined, 10)
	if !ok {
		return nil, fmt.Errorf("cannot parse %q as number", s)
	}

	if result.Sign() < 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	return result, nil
}
