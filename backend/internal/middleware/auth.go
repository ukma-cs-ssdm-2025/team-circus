package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

type userRepository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
}

func AuthMiddleware(userRepo userRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			return
		}

		secretToken := os.Getenv("SECRET_TOKEN")
		if secretToken == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretToken), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			return
		}

		if claims.Subject == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing subject"})
			return
		}

		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid uuid format in token"})
			return
		}

		user, err := userRepo.GetByUUID(c.Request.Context(), uid)
		if err != nil {
			if errors.Is(err, domain.ErrInternal) {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// c.Set("user", user)

		c.Next()
	}
}
