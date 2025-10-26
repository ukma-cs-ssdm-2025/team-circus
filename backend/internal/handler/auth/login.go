package auth

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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth/requests"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
}

// NewLogInHandler handles user login and saves JWT tokens in cookies.
// @Summary User login
// @Description Authenticates user credentials, creates and saves JWT access/refresh tokens in cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.LogInRequest true "User login request"
// @Success 200 {object} map[string]string "Login successful"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/login [post]
func NewLogInHandler(userRepo userRepository, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.LogInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("log in handler: failed to bind request: %v", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		user, err := userRepo.GetByLogin(c, req.Login)
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to log in", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log in"})
			return
		}
		if err != nil {
			logger.Error("failed to log in", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log in"})
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		secretToken := os.Getenv("SECRET_TOKEN")
		if secretToken == "" {
			logger.Error("server misconfiguration")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
			return
		}

		accessExpTime := time.Now().Add(10 * time.Minute)
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   user.UUID.String(),
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
		})

		accessTokenString, err := accessToken.SignedString([]byte(secretToken))
		if err != nil {
			logger.Error("failed to generate access token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
			return
		}

		refreshExpTime := time.Now().Add(30 * 24 * time.Hour)
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   user.UUID.String(),
			ExpiresAt: jwt.NewNumericDate(refreshExpTime),
		})

		refreshTokenString, err := refreshToken.SignedString([]byte(secretToken))
		if err != nil {
			logger.Error("failed to generate refresh token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("accessToken", accessTokenString, int(time.Until(accessExpTime).Seconds()), "/", "", true, true)
		c.SetCookie("refreshToken", refreshTokenString, int(time.Until(refreshExpTime).Seconds()), "/", "", true, true)

		c.JSON(http.StatusOK, gin.H{})
	}
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
	})
}
