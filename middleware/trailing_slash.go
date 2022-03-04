package middleware

import (
	"net/http"
	"strings"
)

// wtf, gorilla ?!
// https://github.com/gorilla/mux/issues/30#issuecomment-43832004
// https://www.husainalshehhi.com/blog/gorilla-mux-trailing-slashes/
func WithTrailingSlash() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			}
			h.ServeHTTP(w, r)
		})
	}
}
