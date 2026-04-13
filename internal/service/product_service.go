package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/client"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/store"
)

var (
	ErrInvalidProductID    = errors.New("invalid product_id")
	ErrProductNotFound     = errors.New("product not found")
	ErrUpstreamUnavailable = errors.New("upstream coinbase service unavailable")
)

type CoinbaseProductGetter interface {
	GetProduct(ctx context.Context, productID string) (model.CoinbaseProduct, error)
}

type ProductService struct {
	coinbaseClient CoinbaseProductGetter
	cache          ProductCache
}

type ProductCache interface {
	GetProduct(ctx context.Context, productID string) (model.ProductResponse, error)
	SetProduct(ctx context.Context, product model.ProductResponse) error
}

func NewProductService(coinbaseClient CoinbaseProductGetter, cache ProductCache) *ProductService {
	return &ProductService{
		coinbaseClient: coinbaseClient,
		cache:          cache,
	}
}

// GetProduct returns a service-owned product response for the requested product ID.
// It uses a read through cache flow: check Redis first, fall back to Coinbase public API on a miss,
// then cache the transformed response before returning it.
// Redis is treated as a cache layer in this version of the service,
// so cache failures are logged but not taking down the endpoint entirely.
func (s *ProductService) GetProduct(ctx context.Context, productID string) (model.ProductResponse, error) {
	normalizedProductID := strings.ToUpper(strings.TrimSpace(productID))
	if normalizedProductID == "" {
		return model.ProductResponse{}, ErrInvalidProductID
	}

	// Prefer the cache first so repeated reads can avoid unnecessary upstream calls.
	if s.cache != nil {
		cachedProduct, err := s.cache.GetProduct(ctx, normalizedProductID)
		switch {
		case err == nil:
			cachedProduct.CacheStatus = "hit"
			return cachedProduct, nil
		case errors.Is(err, store.ErrCacheMiss):
			// Cache misses are expected and should fall through to the Coinbase lookup.
		default:
			// Cache availability should not take down the endpoint in this simplified version.
			log.Printf("redis cache get failed for %s: %v", normalizedProductID, err)
		}
	}

	product, err := s.coinbaseClient.GetProduct(ctx, normalizedProductID)
	if err != nil {
		switch {
		case errors.Is(err, client.ErrProductNotFound):
			return model.ProductResponse{}, fmt.Errorf("%w: %s", ErrProductNotFound, normalizedProductID)
		default:
			return model.ProductResponse{}, fmt.Errorf("%w: %v", ErrUpstreamUnavailable, err)
		}
	}

	// Transform the upstream Coinbase response into the service-owned API contract.
	response := toProductResponse(product)
	response.CacheStatus = "miss"

	// Write back to the cache after a successful upstream fetch.
	if s.cache != nil {
		if err := s.cache.SetProduct(ctx, response); err != nil {
			log.Printf("redis cache set failed for %s: %v", normalizedProductID, err)
		}
	}

	return response, nil
}

func toProductResponse(product model.CoinbaseProduct) model.ProductResponse {
	return model.ProductResponse{
		ProductID:        product.ProductID,
		MarketPair:       buildMarketPair(product.BaseCurrencyID, product.QuoteCurrencyID),
		ProductName:      resolveProductName(product),
		BaseCurrency:     product.BaseCurrencyID,
		QuoteCurrency:    product.QuoteCurrencyID,
		Status:           normalizeStatus(product.Status),
		IsTradingEnabled: isTradingEnabled(product.Status),
		Price:            product.Price,
		PriceChange24H:   product.PricePercentageChange24h,
		RetrievedAt:      time.Now().UTC().Format(time.RFC3339),
		Source:           "coinbase",
	}
}

func buildMarketPair(baseCurrency, quoteCurrency string) string {
	if baseCurrency == "" || quoteCurrency == "" {
		return ""
	}

	return baseCurrency + "/" + quoteCurrency
}

func resolveProductName(product model.CoinbaseProduct) string {
	if product.BaseName != "" {
		return product.BaseName
	}

	if product.DisplayName != "" {
		return product.DisplayName
	}

	return product.ProductID
}

func normalizeStatus(status string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized == "" {
		return "unknown"
	}

	return normalized
}

func isTradingEnabled(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "online":
		return true
	default:
		return false
	}
}
