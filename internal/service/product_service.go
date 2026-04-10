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

	return model.ProductResponse{
		ProductID:                product.ProductID,
		DisplayName:              product.DisplayName,
		BaseCurrencyID:           product.BaseCurrencyID,
		BaseName:                 product.BaseName,
		QuoteCurrencyID:          product.QuoteCurrencyID,
		QuoteName:                product.QuoteName,
		Status:                   product.Status,
		Price:                    product.Price,
		PricePercentageChange24h: product.PricePercentageChange24h,
		BaseIncrement:            product.BaseIncrement,
		QuoteIncrement:           product.QuoteIncrement,
	}, nil
}
