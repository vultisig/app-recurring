package recurring

import (
	"testing"

	rtypes "github.com/vultisig/recipes/types"
	"github.com/vultisig/vultisig-go/common"
)

func TestCreateSendMetaRules_NativeToken(t *testing.T) {
	s := NewSendSpec(nil)

	cfg := map[string]any{
		"asset": map[string]any{
			"chain":   "Ethereum",
			"token":   "native",
			"address": "0x2d63088Dacce3a87b0966982D52141AEe53be224",
		},
		"recipients": []any{
			map[string]any{
				"toAddress": "0x742d35Cc6634C0532925a3b844Bc9e7595f4D215",
				"amount":    "100000000000000",
			},
		},
	}

	rules, err := s.createSendMetaRules(cfg, common.Ethereum, "")
	if err != nil {
		t.Fatalf("createSendMetaRules failed: %v", err)
	}

	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	rule := rules[0]

	// Check resource
	if rule.Resource != "ethereum.send" {
		t.Errorf("expected resource 'ethereum.send', got '%s'", rule.Resource)
	}

	// Find the asset constraint
	var assetConstraint *rtypes.ParameterConstraint
	for _, c := range rule.ParameterConstraints {
		if c.ParameterName == "asset" {
			assetConstraint = c
			break
		}
	}

	if assetConstraint == nil {
		t.Fatal("asset constraint not found")
	}

	// The key fix: for native tokens, FixedValue should be empty string ""
	// This allows the metarule to correctly identify it as a native transfer
	assetValue := assetConstraint.Constraint.GetFixedValue()
	if assetValue != "" {
		t.Errorf("for native token, asset constraint FixedValue should be empty string, got '%s'", assetValue)
	}
}

func TestCreateSendMetaRules_ERC20Token(t *testing.T) {
	s := NewSendSpec(nil)

	tokenAddress := "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" // USDC

	cfg := map[string]any{
		"asset": map[string]any{
			"chain":   "Ethereum",
			"token":   tokenAddress,
			"address": "0x2d63088Dacce3a87b0966982D52141AEe53be224",
		},
		"recipients": []any{
			map[string]any{
				"toAddress": "0x742d35Cc6634C0532925a3b844Bc9e7595f4D215",
				"amount":    "1000000",
			},
		},
	}

	rules, err := s.createSendMetaRules(cfg, common.Ethereum, "")
	if err != nil {
		t.Fatalf("createSendMetaRules failed: %v", err)
	}

	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	rule := rules[0]

	// Find the asset constraint
	var assetConstraint *rtypes.ParameterConstraint
	for _, c := range rule.ParameterConstraints {
		if c.ParameterName == "asset" {
			assetConstraint = c
			break
		}
	}

	if assetConstraint == nil {
		t.Fatal("asset constraint not found")
	}

	// For ERC20 tokens, FixedValue should be the token address
	assetValue := assetConstraint.Constraint.GetFixedValue()
	if assetValue != tokenAddress {
		t.Errorf("for ERC20 token, asset constraint FixedValue should be '%s', got '%s'", tokenAddress, assetValue)
	}
}
