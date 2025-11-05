package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/httpx"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/responses"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
// @Router /users/{uuid} [get]
func NewGetUserHandler(service getUserService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"get user handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
			return
		}

		user, err := service.GetByUUID(c.Request.Context(), userUUID)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to get user",
				LogMessage: "failed to get user",
				LogLevel:   zapcore.ErrorLevel,
			},
			httpx.RequestContextFields(c, zap.String("user_uuid", userUUID.String())),
			httpx.ResponseSpec{
				Target:     domain.ErrUserNotFound,
				Status:     http.StatusNotFound,
				Message:    "user not found",
				LogMessage: "user not found",
				LogLevel:   zapcore.WarnLevel,
			},
		) {
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
// @Router /users [get]
func NewGetAllUsersHandler(service getAllUsersService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := service.GetAll(c.Request.Context())
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to get users",
				LogMessage: "failed to get users",
				LogLevel:   zapcore.ErrorLevel,
			},
			httpx.RequestContextFields(c),
		) {
			return
		}

		response := responses.GetAllUsersResponse{
			Users: mapUsersToGetAllResponse(users),
		}

		c.JSON(http.StatusOK, response)
	}
}
