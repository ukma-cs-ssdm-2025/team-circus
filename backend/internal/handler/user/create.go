package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/requests"
)

type createUserService interface {
	Create(ctx context.Context, login string, email string, password string) (*domain.User, error)
}

// NewCreateUserHandler creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided login, email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "User creation request"
// @Success 201 {object} responses.CreateUserResponse "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation failed"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/users [post]
func NewCreateUserHandler(service createUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("create user handler: failed to bind request: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("create user handler: validation failed: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		user, err := service.Create(c, req.Login, req.Email, req.Password)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		response := mapUserToCreateResponse(user)

		c.JSON(http.StatusCreated, response)
	}
}
