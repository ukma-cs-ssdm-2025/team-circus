package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/requests"
	"go.uber.org/zap"
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
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("update user handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
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

		user, err := service.Update(c.Request.Context(), parsedUUID, req.Login, req.Email, req.Password)
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Warn("user not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to update user", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}
		if err != nil {
			logger.Error("failed to update user", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		response := mapUserToUpdateResponse(user)

		c.JSON(http.StatusOK, response)
	}
}
