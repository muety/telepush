package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}
	return m[0](Chain(f, m[1:cap(m)]...))
}

func CheckMethod(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(415)
			return
		}
		next(w, r)
	}
}
