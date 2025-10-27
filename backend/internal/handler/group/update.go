package group

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/requests"
	"go.uber.org/zap"
)

type updateGroupService interface {
	Update(ctx context.Context, userUUID, groupUUID uuid.UUID, name string) (*domain.Group, error)
}

// NewUpdateGroupHandler updates a group by UUID
// @Summary Update a group by UUID
// @Description Update a specific group's name by its UUID
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Param request body requests.UpdateGroupRequest true "Group update request"
// @Success 200 {object} responses.UpdateGroupResponse "Group updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format or validation failed"
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid} [put]
func NewUpdateGroupHandler(service updateGroupService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUIDValue, exists := c.Get("user_uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user context missing"})
			return
		}

		userUUID, ok := userUUIDValue.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
			return
		}

		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("update group handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		var req requests.UpdateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("update group handler: failed to bind request: %v", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("update group handler: validation failed: %v", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		group, err := service.Update(c, userUUID, parsedUUID, req.Name)
		if errors.Is(err, domain.ErrGroupNotFound) {
			logger.Warn("group not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to update group", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update group"})
			return
		}
		if err != nil {
			logger.Error("failed to update group", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update group"})
			return
		}

		response := mapGroupToUpdateResponse(group)

		c.JSON(http.StatusOK, response)
	}
}
