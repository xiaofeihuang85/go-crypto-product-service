package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/client"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
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
}

func NewProductService(coinbaseClient CoinbaseProductGetter) *ProductService {
	return &ProductService{
		coinbaseClient: coinbaseClient,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, productID string) (model.ProductResponse, error) {
	normalizedProductID := strings.ToUpper(strings.TrimSpace(productID))
	if normalizedProductID == "" {
		return model.ProductResponse{}, ErrInvalidProductID
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

	return toProductResponse(product), nil
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
