package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"go.uber.org/zap"
)

// NewUpdateRefreshTokenHandler handles refresh token requests
// @Summary Refresh access token
// @Description Validates the refresh token cookie and issues a new access/refresh token pair
// @Tags auth
// @Accept */*
// @Produce json
// @Param refreshToken cookie string true "Refresh token"
// @Success 200 {object} map[string]string "Tokens refreshed successfully"
// @Failure 401 {object} map[string]string "Invalid or expired refresh token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/refresh [post]
func NewRefreshTokenHandler(userRepo userRepository, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var req requests.RefreshTokenRequest
		// if err := c.ShouldBindJSON(&req); err != nil {
		// 	err = fmt.Errorf("refresh token handler: failed to bind request: %v", err)
		// 	logger.Error("failed to bind request", zap.Error(err))
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		// 	return
		// }

		tokenString, err := c.Cookie("refreshToken")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token required"})
			return
		}

		secretToken := os.Getenv("SECRET_TOKEN")
		if secretToken == "" {
			logger.Error("server misconfiguration")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretToken), nil
		})
		if err != nil {
			logger.Error("invalid refresh token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			logger.Error("invalid refresh token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token claims"})
			return
		}

		if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			logger.Error("refresh token expired")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
			return
		}

		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			logger.Error("invalid uuid format in refresh token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid uuid format in refresh token"})
			return
		}

		user, err := userRepo.GetByUUID(c.Request.Context(), uid)
		if err != nil {
			if errors.Is(err, domain.ErrInternal) {
				logger.Error("failed to fetch user from refresh token", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user from refresh token"})
				return
			}
			logger.Error("invalid refresh token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		accessTokenExpTime := time.Now().Add(10 * time.Minute)
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   user.UUID.String(),
			ExpiresAt: jwt.NewNumericDate(accessTokenExpTime),
		})

		accessTokenString, err := accessToken.SignedString([]byte(secretToken))
		if err != nil {
			logger.Error("failed to generate access token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
			return
		}

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   user.UUID.String(),
			ExpiresAt: claims.ExpiresAt,
		})

		refreshTokenString, err := refreshToken.SignedString([]byte(secretToken))
		if err != nil {
			logger.Error("failed to generate refresh token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("accessToken", accessTokenString, int(time.Until(accessTokenExpTime).Seconds()), "/", "", true, true)
		c.SetCookie("refreshToken", refreshTokenString, int(time.Until(claims.ExpiresAt.Time).Seconds()), "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{})
	}
}
