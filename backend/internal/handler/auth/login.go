package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth/requests"
)

type userRepository interface {
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
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
// @Router /api/v1/auth/login [post]
func NewLogInHandler(userRepo userRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.LogInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("log in handler: failed to bind request: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		user, err := userRepo.GetByLogin(c, req.Login)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log in"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log in"})
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if req.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		secretToken := os.Getenv("SECRET_TOKEN")
		if secretToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
			return
		}

		accessExpirationTime := time.Now().Add(10 * time.Minute)
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   user.UUID.String(),
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		})

		accessTokenString, err := accessToken.SignedString([]byte(secretToken))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
			return
		}

		refreshExpirationTime := time.Now().Add(30 * 24 * time.Hour)
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   user.UUID.String(),
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		})

		refreshTokenString, err := refreshToken.SignedString([]byte(secretToken))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("accessToken", accessTokenString, 60*10, "", "", true, true)
		c.SetCookie("refreshToken", refreshTokenString, 60*60*24*30, "", "", true, true)

		c.JSON(http.StatusOK, gin.H{})
	}
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
	})
}
