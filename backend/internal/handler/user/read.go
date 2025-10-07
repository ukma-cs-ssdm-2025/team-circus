package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/responses"
)

type getUserService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
}

type getAllUsersService interface {
	GetAll(ctx context.Context) ([]*domain.User, error)
}

// NewGetUserHandler retrieves a user by UUID
// @Summary Get a user by UUID
// @Description Retrieve a specific user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Success 200 {object} responses.GetUserResponse "User retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/users/{uuid} [get]
func NewGetUserHandler(service getUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("get user handler: failed to parse uuid: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		user, err := service.GetByUUID(c, parsedUUID)
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}

		response := mapUserToGetResponse(user)

		c.JSON(http.StatusOK, response)
	}
}

// NewGetAllUsersHandler retrieves all users
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllUsersResponse "Users retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/users [get]
func NewGetAllUsersHandler(service getAllUsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := service.GetAll(c)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
			return
		}

		response := responses.GetAllUsersResponse{
			Users: mapUsersToGetAllResponse(users),
		}

		c.JSON(http.StatusOK, response)
	}
}
