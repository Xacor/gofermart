package jwt

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

func (ctxKey) String() string {
	return "ctxKey"
}

const UserIDKey ctxKey = "userID"

func WithJWTAuth(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, tokenString, ok := strings.Cut(r.Header.Get("Authorization"), "Bearer ")
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claims, err := ValidateToken(tokenString, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromCtx(ctx context.Context) int {
	return ctx.Value(UserIDKey).(int)
}
