package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/service"
)

type stubProductService struct {
	product       model.ProductResponse
	err           error
	lastProductID string
}

func (s *stubProductService) GetProduct(_ context.Context, productID string) (model.ProductResponse, error) {
	s.lastProductID = productID
	return s.product, s.err
}

func TestProductHandlerReturnsBadRequestForInvalidPath(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/products/", nil)
	rec := httptest.NewRecorder()

	productHandler(&stubProductService{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	var response model.APIErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != "invalid_product_id" {
		t.Fatalf("expected invalid_product_id, got %q", response.Code)
	}
}

func TestProductHandlerReturnsBadGatewayForUpstreamFailure(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/products/BTC-USD", nil)
	rec := httptest.NewRecorder()
	productService := &stubProductService{err: service.ErrUpstreamUnavailable}

	productHandler(productService).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.Code)
	}

	var response model.APIErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Code != "upstream_unavailable" {
		t.Fatalf("expected upstream_unavailable, got %q", response.Code)
	}

	if response.Path != "/products/BTC-USD" {
		t.Fatalf("expected request path to be echoed, got %q", response.Path)
	}
}

func TestProductHandlerReturnsProductResponse(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/products/BTC-USD", nil)
	rec := httptest.NewRecorder()
	productService := &stubProductService{
		product: model.ProductResponse{
			ProductID:   "BTC-USD",
			MarketPair:  "BTC/USD",
			CacheStatus: "hit",
			Source:      "coinbase",
		},
	}

	productHandler(productService).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if productService.lastProductID != "BTC-USD" {
		t.Fatalf("expected service to receive BTC-USD, got %q", productService.lastProductID)
	}

	var response model.ProductResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ProductID != "BTC-USD" {
		t.Fatalf("expected BTC-USD, got %q", response.ProductID)
	}
}

func TestProductHandlerReturnsNotFoundForMissingProduct(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/products/DOES-NOT-EXIST", nil)
	rec := httptest.NewRecorder()
	productService := &stubProductService{err: fmt.Errorf("%w: DOES-NOT-EXIST", service.ErrProductNotFound)}

	productHandler(productService).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
