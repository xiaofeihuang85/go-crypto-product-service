package api

import "net/http"

func NewRouter(serviceName string, productService ProductService) http.Handler {
	mux := http.NewServeMux()
	uiHandler := newUIHandler()

	mux.HandleFunc("/", rootHandler(serviceName))
	mux.HandleFunc("/health", healthHandler(serviceName))
	mux.HandleFunc("/products/", productHandler(productService))
	mux.Handle("/app", uiHandler)
	mux.Handle("/app/", uiHandler)

	return mux
}
