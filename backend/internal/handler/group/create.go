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

type createGroupService interface {
	Create(ctx context.Context, ownerUUID uuid.UUID, name string) (*domain.Group, error)
}

// NewCreateGroupHandler creates a new group
// @Summary Create a new group
// @Description Create a new group with the provided name
// @Tags groups
// @Accept json
// @Produce json
// @Param request body requests.CreateGroupRequest true "Group creation request"
// @Success 201 {object} responses.CreateGroupResponse "Group created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation failed"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups [post]
func NewCreateGroupHandler(service createGroupService, logger *zap.Logger) gin.HandlerFunc {
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

		var req requests.CreateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("create group handler: failed to bind request: %v", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("create group handler: validation failed: %v", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		group, err := service.Create(c, userUUID, req.Name)
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to create group", zap.Error(err), zap.String("name", req.Name))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
			return
		}
		if err != nil {
			logger.Error("failed to create group", zap.Error(err), zap.String("name", req.Name))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
			return
		}

		response := mapGroupToCreateResponse(group)

		c.JSON(http.StatusCreated, response)
	}
}
