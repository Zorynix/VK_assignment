package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"vk.com/m/auth"
)

func AuthMiddleware(next http.Handler, allowedRoles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeaderParts[1]
		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Error().Err(err).Str("token", tokenStr).Msg("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		roleIsAllowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				roleIsAllowed = true
				break
			}
		}

		if !roleIsAllowed {
			log.Warn().Str("role", claims.Role).Msg("Attempt to access with insufficient permissions")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else {
			log.Info().Str("role", claims.Role).Msg("Access granted")
		}

		next.ServeHTTP(w, r)
	})
}
