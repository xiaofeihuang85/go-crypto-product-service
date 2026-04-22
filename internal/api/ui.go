package api

import "net/http"

func newUIHandler() http.Handler {
	fileServer := http.FileServer(http.Dir("web"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/app", "/app/":
			http.ServeFile(w, r, "web/index.html")
			return
		default:
			if len(r.URL.Path) > len("/app/") {
				http.StripPrefix("/app/", fileServer).ServeHTTP(w, r)
				return
			}
			http.NotFound(w, r)
		}
	})
}
