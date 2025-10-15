package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth/requests"
)

// NewUpdateRefreshTokenHandler handles refresh token requests
// @Summary Refresh access token
// @Description Validates a refresh token and issues a new access and refresh token pair
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} map[string]string "Tokens refreshed successfully"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/refresh [post]
func NewRefreshTokenHandler(userRepo userRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("refresh token handler: failed to bind request: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		secretToken := os.Getenv("SECRET_TOKEN")
		if secretToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
			return
		}

		token, err := jwt.ParseWithClaims(req.RefreshToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretToken), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token claims"})
			return
		}

		if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
			return
		}

		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid uuid format in access token"})
			return
		}

		accessTokenExpTime := time.Now().Add(10 * time.Minute)
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   uid.String(),
			ExpiresAt: jwt.NewNumericDate(accessTokenExpTime),
		})

		accessTokenString, err := accessToken.SignedString([]byte(secretToken))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
			return
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   uid.String(),
			ExpiresAt: claims.ExpiresAt,
		})

		refreshTokenString, err := refreshToken.SignedString([]byte(secretToken))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessTokenString,
			"refresh_token": refreshTokenString,
		})
	}
}
