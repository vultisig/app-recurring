package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type swapService struct {
	providers []Provider
}

func newSwapService(providers []Provider) *swapService {
	return &swapService{
		providers: providers,
	}
}

type QuoteWinner struct {
	Provider Provider
	Quote    *QuoteResult
}

func (s *swapService) FindBestQuote(
	ctx context.Context,
	from From,
	to To,
) (*QuoteWinner, error) {
	if len(s.providers) == 0 {
		return nil, fmt.Errorf("no providers available")
	}

	providers := s.providers

	type quoteResult struct {
		quote    *QuoteResult
		provider Provider
		err      error
		name     string
	}
	results := make([]quoteResult, len(providers))

	g, ctx := errgroup.WithContext(ctx)

	for _i, _provider := range providers {
		i, provider := _i, _provider
		g.Go(func() error {
			q, err := provider.Quote(ctx, from, to)

			results[i] = quoteResult{
				quote:    q,
				provider: provider,
				err:      err,
				name:     provider.Name(),
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("errgroup failed: %w", err)
	}

	var bestQuote *QuoteResult
	var bestProvider Provider
	var bestAmountOut = big.NewInt(0)
	var bestProviderName string
	var lastErr error

	for _, result := range results {
		if result.err != nil {
			logrus.WithFields(logrus.Fields{
				"provider": result.name,
				"error":    result.err.Error(),
			}).Debug("evm swap quote failed")
			lastErr = result.err
			continue
		}

		logrus.WithFields(logrus.Fields{
			"provider":  result.name,
			"amountOut": result.quote.AmountOut.String(),
		}).Debug("evm swap quote succeeded")

		if result.quote.AmountOut != nil && result.quote.AmountOut.Cmp(bestAmountOut) > 0 {
			bestAmountOut = result.quote.AmountOut
			bestQuote = result.quote
			bestProvider = result.provider
			bestProviderName = result.name
		}
	}

	if bestQuote != nil {
		logrus.WithFields(logrus.Fields{
			"selectedProvider": bestProviderName,
			"bestAmountOut":    bestAmountOut.String(),
			"spender":          bestQuote.Spender.Hex(),
		}).Info("selected best swap quote")
		return &QuoteWinner{
			Provider: bestProvider,
			Quote:    bestQuote,
		}, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all providers failed, last error: %w", lastErr)
	}
	return nil, fmt.Errorf("no valid quotes found")
}

