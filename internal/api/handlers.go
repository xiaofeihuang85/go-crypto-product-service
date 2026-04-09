package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/xiaofeihuang85/go-crypto-product-service/internal/model"
)

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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(v)
}
