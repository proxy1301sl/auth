package middleware

import (
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.Header.Get("Authorization")
		_, token, ok := strings.Cut(req, "Bearer ")
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

	})
}
