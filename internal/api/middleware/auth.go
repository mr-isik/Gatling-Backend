package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mr-isik/gatling-backend/internal/api/httputil"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/service"
)

type contextKey string

const claimsKey contextKey = "claims"

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputil.JSONError(w, http.StatusUnauthorized, domain.ErrUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httputil.JSONError(w, http.StatusUnauthorized, domain.ErrUnauthorized)
				return
			}

			tokenStr := parts[1]
			claims, err := service.ValidateToken(tokenStr, jwtSecret)
			if err != nil {
				httputil.JSONError(w, http.StatusUnauthorized, domain.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaimsFromContext(ctx context.Context) *service.Claims {
	claims, ok := ctx.Value(claimsKey).(*service.Claims)
	if !ok {
		return nil
	}
	return claims
}

func GetUserIDFromContext(ctx context.Context) string {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		return ""
	}
	return claims.UserID
}
