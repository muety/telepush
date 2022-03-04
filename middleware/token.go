package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func WithToken(srcKey, dstKey string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := mux.Vars(r)[srcKey]
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("missing recipient token"))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, dstKey, token)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
