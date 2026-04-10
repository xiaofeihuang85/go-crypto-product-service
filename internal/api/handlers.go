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
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "resource not found",
			})
			return
		}

		writeJSON(w, http.StatusOK, model.ServiceInfoResponse{
			Service: serviceName,
			Status:  "ready",
			Message: "phase 2 http bootstrap is running",
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
			writeJSON(w, http.StatusMethodNotAllowed, model.APIErrorResponse{
				Error: "method not allowed",
			})
			return
		}

		productID, ok := parseProductID(r.URL.Path)
		if !ok {
			writeJSON(w, http.StatusBadRequest, model.APIErrorResponse{
				Error: "product_id must be provided as /products/{product_id}",
			})
			return
		}

		product, err := productService.GetProduct(r.Context(), productID)
		if err != nil {
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, service.ErrInvalidProductID):
				status = http.StatusBadRequest
			case errors.Is(err, service.ErrProductNotFound):
				status = http.StatusNotFound
			case errors.Is(err, service.ErrUpstreamUnavailable):
				status = http.StatusBadGateway
			}

			writeJSON(w, status, model.APIErrorResponse{
				Error: err.Error(),
			})
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
