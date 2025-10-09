package enhanced_router

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

// Claims represents the JWT claims.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT for a given username.
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// AuthMiddleware creates a middleware for JWT authentication.
func AuthMiddleware() MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			tokenString, ok := ctx.Value("token").(string)
			if !ok {
				return nil, fmt.Errorf("token not found")
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					return nil, fmt.Errorf("invalid token signature")
				}
				return nil, fmt.Errorf("bad token")
			}

			if !token.Valid {
				return nil, fmt.Errorf("invalid token")
			}

			// Add username to context for downstream handlers
			ctx = context.WithValue(ctx, "username", claims.Username)

			return next(ctx, req)
		}
	}
}
