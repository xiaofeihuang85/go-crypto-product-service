package api

import "net/http"

func NewRouter(serviceName string, productService ProductService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler(serviceName))
	mux.HandleFunc("/health", healthHandler(serviceName))
	mux.HandleFunc("/products/", productHandler(productService))

	return mux
}
