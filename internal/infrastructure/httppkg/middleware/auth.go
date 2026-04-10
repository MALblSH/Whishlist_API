package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MALblSH/Wishlist_API/internal/infrastructure/tokens"
)

type contextKeyType string

const (
	claimsKey            contextKeyType = "claims"
	authHeaderPartsCount                = 2
	authSchemeIndex                     = 0
	authTokenIndex                      = 1
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractTokenFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), claimsKey, tokenString)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", tokens.ErrTokenMissing
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != authHeaderPartsCount || !strings.EqualFold(parts[authSchemeIndex], "Bearer") {
		return "", tokens.ErrTokenInvalid
	}
	tokenString := parts[authTokenIndex]
	if tokenString == "" {
		return "", tokens.ErrTokenMissing
	}
	return tokenString, nil
}

func ClaimsFromContext(ctx context.Context) (string, bool) {
	claims, ok := ctx.Value(claimsKey).(string)
	return claims, ok
}
