package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/httpx"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/requests"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type updateUserService interface {
	Update(ctx context.Context, uuid uuid.UUID, login string, email string, password string) (*domain.User, error)
}

// NewUpdateUserHandler updates a user by UUID
// @Summary Update a user by UUID
// @Description Update a specific user's login, email and password by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Param request body requests.UpdateUserRequest true "User update request"
// @Success 200 {object} responses.UpdateUserResponse "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format or validation failed"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{uuid} [put]
func NewUpdateUserHandler(service updateUserService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"update user handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
			return
		}

		var req requests.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("update user handler: failed to bind request: %v", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("update user handler: validation failed: %v", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		user, err := service.Update(c.Request.Context(), userUUID, req.Login, req.Email, req.Password)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to update user",
				LogMessage: "failed to update user",
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

		response := mapUserToUpdateResponse(user)

		c.JSON(http.StatusOK, response)
	}
}
