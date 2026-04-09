package api

import "net/http"

func NewRouter(serviceName string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler(serviceName))
	mux.HandleFunc("/health", healthHandler(serviceName))

	return mux
}
