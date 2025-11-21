package group

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/responses"
	"go.uber.org/zap"
)

type deleteGroupService interface {
	Delete(ctx context.Context, uuid uuid.UUID) error
}

// NewDeleteGroupHandler deletes a group by UUID
// @Summary Delete a group by UUID
// @Description Delete a specific group by its UUID
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Success 200 {object} responses.DeleteGroupResponse "Group deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid} [delete]
func NewDeleteGroupHandler(service deleteGroupService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("delete group handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		err = service.Delete(c.Request.Context(), parsedUUID)
		if errors.Is(err, domain.ErrGroupNotFound) {
			logger.Warn("group not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to delete group", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group"})
			return
		}
		if err != nil {
			logger.Error("failed to delete group", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group"})
			return
		}

		response := responses.DeleteGroupResponse{
			Message: "Group deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
