package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/responses"
	userrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
	"go.uber.org/zap"
)

type getUserService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
}

type getAllUsersService interface {
	GetAll(ctx context.Context, params userrepo.PageParams) ([]*domain.User, error)
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
// @Router /users/{uuid} [get]
func NewGetUserHandler(service getUserService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("get user handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		user, err := service.GetByUUID(c, parsedUUID)
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Warn("user not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to get user", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}
		if err != nil {
			logger.Error("failed to get user", zap.Error(err), zap.String("uuid", uuidParam))
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
// @Param limit query int false "Max users per page (1-100, defaults to 50)"
// @Param offset query int false "Offset to start listing from (defaults to 0)"
// @Success 200 {object} responses.GetAllUsersResponse "Users retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Failure 400 {object} map[string]interface{} "Invalid pagination parameters"
// @Router /users [get]
func NewGetAllUsersHandler(service getAllUsersService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageParams, err := parsePageParams(c)
		if err != nil {
			logger.Warn("invalid pagination params", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		users, err := service.GetAll(c, pageParams)
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to get users", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
			return
		}
		if err != nil {
			logger.Error("failed to get users", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
			return
		}

		response := responses.GetAllUsersResponse{
			Users: mapUsersToGetAllResponse(users),
			Meta: responses.PageMeta{
				Limit:  pageParams.Limit,
				Offset: pageParams.Offset,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func parsePageParams(c *gin.Context) (userrepo.PageParams, error) {
	params := userrepo.PageParams{
		Limit:  userrepo.DefaultPageLimit,
		Offset: 0,
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return params, errors.New("invalid limit value")
		}
		params.Limit = limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return params, errors.New("invalid offset value")
		}
		params.Offset = offset
	}

	return params.Normalize(), nil
}
