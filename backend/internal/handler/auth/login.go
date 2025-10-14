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

// NewLogInHandler handles user login and saves JWT token
// @Summary User login
// @Description Authenticate user, create and save JWT token
// @Tags login
// @Accept json
// @Produce json
// @Param request body requests.LogInRequest true "User login request"
// @Success 200 {object} {} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation failed"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/login [post]
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

		expirationTime := time.Now().Add(15 * time.Minute)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.UUID,
			"exp": expirationTime.Unix(),
		})

		tokenString, err := token.SignedString([]byte(secretToken))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 60*15, "", "", true, true)

		c.JSON(http.StatusOK, gin.H{})
	}
}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
	})
}
