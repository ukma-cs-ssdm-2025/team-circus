package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/responses"
	"go.uber.org/zap"
)

type deleteUserService interface {
	Delete(ctx context.Context, uuid uuid.UUID) error
}

// NewDeleteUserHandler deletes a user by UUID
// @Summary Delete a user by UUID
// @Description Delete a specific user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "User UUID"
// @Success 200 {object} responses.DeleteUserResponse "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{uuid} [delete]
func NewDeleteUserHandler(service deleteUserService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("delete user handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		err = service.Delete(c.Request.Context(), parsedUUID)
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.Warn("user not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to delete user", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			return
		}
		if err != nil {
			logger.Error("failed to delete user", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			return
		}

		response := responses.DeleteUserResponse{
			Message: "User deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
