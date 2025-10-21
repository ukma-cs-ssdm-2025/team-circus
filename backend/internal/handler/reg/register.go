package reg

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg/requests"
	"go.uber.org/zap"
)

type regService interface {
	Register(ctx context.Context, login string, email string, password string) (*domain.User, error)
}

// NewCreateUserHandler registers a new user
// @Summary Register a new user
// @Description Register a new user with the provided login, email and password
// @Tags registration
// @Accept json
// @Produce json
// @Param request body requests.RegRequest true "User registration request"
// @Success 201 {object} responses.RegResponse "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation failed"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /signup [post]
func NewRegHandler(service regService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.RegRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("register handler: failed to bind request: %v", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("register handler: validation failed: %v", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		user, err := service.Register(c, req.Login, req.Email, req.Password)
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to register",
				zap.Error(err),
				zap.String("login", req.Login),
				zap.String("email", req.Email),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
			return
		}
		if err != nil {
			logger.Error("failed to register",
				zap.Error(err),
				zap.String("login", req.Login),
				zap.String("email", req.Email),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
			return
		}

		response := mapUserToRegResponse(user)

		c.JSON(http.StatusCreated, response)
	}
}
