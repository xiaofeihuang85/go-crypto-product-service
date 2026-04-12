package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
	"github.com/xiaofeihuang85/go-crypto-product-service/internal/service"
)

type ProductService interface {
	GetProduct(ctx context.Context, productID string) (model.ProductResponse, error)
}

func rootHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			writeAPIError(w, http.StatusNotFound, "resource_not_found", "resource not found", r.URL.Path, "")
			return
		}

		writeJSON(w, http.StatusOK, model.ServiceInfoResponse{
			Service: serviceName,
			Status:  "ready",
			Message: "phase 6 product cache endpoint is running",
		})
	}
}

func healthHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, model.HealthResponse{
			Status:  "ok",
			Service: serviceName,
			Time:    time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func productHandler(productService ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", r.URL.Path, "only GET is supported")
			return
		}

		productID, ok := parseProductID(r.URL.Path)
		if !ok {
			writeAPIError(w, http.StatusBadRequest, "invalid_product_id", "product_id must be provided as /products/{product_id}", r.URL.Path, "")
			return
		}

		product, err := productService.GetProduct(r.Context(), productID)
		if err != nil {
			status := http.StatusInternalServerError
			code := "internal_error"
			details := ""
			switch {
			case errors.Is(err, service.ErrInvalidProductID):
				status = http.StatusBadRequest
				code = "invalid_product_id"
			case errors.Is(err, service.ErrProductNotFound):
				status = http.StatusNotFound
				code = "product_not_found"
			case errors.Is(err, service.ErrUpstreamUnavailable):
				status = http.StatusBadGateway
				code = "upstream_unavailable"
				details = "coinbase lookup failed after cache lookup did not return a value"
			}

			writeAPIError(w, status, code, err.Error(), r.URL.Path, details)
			return
		}

		writeJSON(w, http.StatusOK, product)
	}
}

func parseProductID(path string) (string, bool) {
	const prefix = "/products/"

	if !strings.HasPrefix(path, prefix) {
		return "", false
	}

	productID := strings.TrimSpace(strings.TrimPrefix(path, prefix))
	if productID == "" || strings.Contains(productID, "/") {
		return "", false
	}

	return productID, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(v)
}

func writeAPIError(w http.ResponseWriter, status int, code, message, path, details string) {
	writeJSON(w, status, model.APIErrorResponse{
		Code:    code,
		Error:   message,
		Path:    path,
		Details: details,
	})
}
