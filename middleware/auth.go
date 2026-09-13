package middleware

import (
	"context"
	"net/http"
	"strings"

	"https/github.com/proxy1301sl/auth/token"
)

type ContextKey string

const UserIdKey ContextKey = "userId"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.Header.Get("Authorization")
		tokenf := strings.TrimPrefix(req, "Bearer ")
		if tokenf == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if tokenf == "" || tokenf != req {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims, err := token.VerifyJWT(tokenf)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIdKey, claims.Subject)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
