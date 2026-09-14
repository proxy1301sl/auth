package middlewareauth

import (
	"net/http"
	"strings"

	authctx "https/github.com/proxy1301sl/auth"
	"https/github.com/proxy1301sl/auth/token"
)

type ContextKey string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.Header.Get("Authorization")
		if !strings.HasPrefix(req, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenf := strings.TrimPrefix(req, "Bearer ")
		claims, err := token.VerifyJWT(tokenf)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := authctx.WithUserId(r.Context(), claims.Subject)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
