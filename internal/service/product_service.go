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

func (s *ProductService) GetProduct(ctx context.Context, productID string) (model.ProductResponse, error) {
	normalizedProductID := strings.ToUpper(strings.TrimSpace(productID))
	if normalizedProductID == "" {
		return model.ProductResponse{}, ErrInvalidProductID
	}

	if s.cache != nil {
		cachedProduct, err := s.cache.GetProduct(ctx, normalizedProductID)
		switch {
		case err == nil:
			cachedProduct.CacheStatus = "hit"
			return cachedProduct, nil
		case errors.Is(err, store.ErrCacheMiss):
		default:
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

	response := toProductResponse(product)
	response.CacheStatus = "miss"

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
