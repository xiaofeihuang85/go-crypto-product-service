package service

import (
	"context"
	"errors"
	"testing"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/client"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/store"
)

type stubCoinbaseClient struct {
	product    model.CoinbaseProduct
	err        error
	lastLookup string
}

func (s *stubCoinbaseClient) GetProduct(_ context.Context, productID string) (model.CoinbaseProduct, error) {
	s.lastLookup = productID
	return s.product, s.err
}

type stubProductCache struct {
	product       model.ProductResponse
	getErr        error
	setErr        error
	lastSet       model.ProductResponse
	setCalled     bool
	lastRequested string
}

func (s *stubProductCache) GetProduct(_ context.Context, productID string) (model.ProductResponse, error) {
	s.lastRequested = productID
	return s.product, s.getErr
}

func (s *stubProductCache) SetProduct(_ context.Context, product model.ProductResponse) error {
	s.setCalled = true
	s.lastSet = product
	return s.setErr
}

func TestGetProductReturnsCachedResponseOnHit(t *testing.T) {
	t.Parallel()

	cache := &stubProductCache{
		product: model.ProductResponse{
			ProductID:   "BTC-USD",
			CacheStatus: "miss",
			Source:      "coinbase",
		},
	}
	coinbaseClient := &stubCoinbaseClient{}
	service := NewProductService(coinbaseClient, cache)

	product, err := service.GetProduct(context.Background(), "btc-usd")
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}

	if product.CacheStatus != "hit" {
		t.Fatalf("expected cache hit, got %q", product.CacheStatus)
	}

	if coinbaseClient.lastLookup != "" {
		t.Fatalf("expected Coinbase client to be skipped on cache hit, got lookup for %q", coinbaseClient.lastLookup)
	}
}

func TestGetProductFetchesAndCachesOnMiss(t *testing.T) {
	t.Parallel()

	cache := &stubProductCache{getErr: store.ErrCacheMiss}
	coinbaseClient := &stubCoinbaseClient{
		product: model.CoinbaseProduct{
			ProductID:                "BTC-USD",
			BaseCurrencyID:           "BTC",
			QuoteCurrencyID:          "USD",
			BaseName:                 "Bitcoin",
			Status:                   "online",
			Price:                    "70000.01",
			PricePercentageChange24h: "-1.25",
		},
	}
	service := NewProductService(coinbaseClient, cache)

	product, err := service.GetProduct(context.Background(), "btc-usd")
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}

	if coinbaseClient.lastLookup != "BTC-USD" {
		t.Fatalf("expected normalized lookup BTC-USD, got %q", coinbaseClient.lastLookup)
	}

	if product.MarketPair != "BTC/USD" {
		t.Fatalf("expected market pair BTC/USD, got %q", product.MarketPair)
	}

	if product.ProductName != "Bitcoin" {
		t.Fatalf("expected product name Bitcoin, got %q", product.ProductName)
	}

	if product.CacheStatus != "miss" {
		t.Fatalf("expected cache miss, got %q", product.CacheStatus)
	}

	if product.Source != "coinbase" {
		t.Fatalf("expected source coinbase, got %q", product.Source)
	}

	if product.RetrievedAt == "" {
		t.Fatal("expected retrieved_at to be populated")
	}

	if !cache.setCalled {
		t.Fatal("expected cache set to be called after upstream fetch")
	}

	if cache.lastSet.CacheStatus != "miss" {
		t.Fatalf("expected cached object to record miss on initial fetch, got %q", cache.lastSet.CacheStatus)
	}
}

func TestGetProductReturnsNotFoundForMissingCoinbaseProduct(t *testing.T) {
	t.Parallel()

	service := NewProductService(&stubCoinbaseClient{err: client.ErrProductNotFound}, &stubProductCache{getErr: store.ErrCacheMiss})

	_, err := service.GetProduct(context.Background(), "DOES-NOT-EXIST")
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestGetProductReturnsInvalidProductIDForBlankInput(t *testing.T) {
	t.Parallel()

	service := NewProductService(&stubCoinbaseClient{}, nil)

	_, err := service.GetProduct(context.Background(), "   ")
	if !errors.Is(err, ErrInvalidProductID) {
		t.Fatalf("expected ErrInvalidProductID, got %v", err)
	}
}
