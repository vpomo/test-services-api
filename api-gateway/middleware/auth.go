package middleware

import (
	"context"
	"net/http"

	"main/internal/proto/user"
)

func AuthMiddleware(userClient user.UserServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Authorization token required", http.StatusUnauthorized)
				return
			}

			resp, err := userClient.ValidateToken(r.Context(), &user.ValidateTokenRequest{
				Token: token,
			})
			if err != nil || !resp.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
