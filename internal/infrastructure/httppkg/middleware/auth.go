package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/MALblSH/Wishlist_API/internal/infrastructure/tokens"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/tokens/access"
)

type contextKeyType string

const (
	claimsKey            contextKeyType = "claims"
	authHeaderPartsCount int            = 2
	authSchemeIndex      int            = 0
	authTokenIndex       int            = 1
)

func AuthMiddleWare(manager *access.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractTokenFromRequest(r)
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}
			var claims *access.Claims
			claims, err = manager.ParseAccessToken(tokenString)
			if err != nil {
				if errors.Is(err, tokens.ErrTokenExpired) {
					http.Error(w, "token expired", http.StatusUnauthorized)
					return
				}
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", tokens.ErrTokenMissing
	}
	parts := strings.SplitN(authHeader, " ", authHeaderPartsCount)
	if len(parts) != authHeaderPartsCount || !strings.EqualFold(parts[authSchemeIndex], "Bearer") {
		return "", tokens.ErrTokenInvalid
	}
	tokenString := parts[authTokenIndex]
	if tokenString == "" {
		return "", tokens.ErrTokenMissing
	}
	return tokenString, nil
}

func ClaimsFromContext(ctx context.Context) (*access.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*access.Claims)
	return claims, ok
}
