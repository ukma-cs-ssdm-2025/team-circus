package group

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/httpx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		groupUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"update group handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
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

		group, err := service.Update(c.Request.Context(), userUUID, groupUUID, req.Name)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to update group",
				LogMessage: "failed to update group",
				LogLevel:   zapcore.ErrorLevel,
			},
			httpx.RequestContextFields(c, zap.String("group_uuid", groupUUID.String())),
			httpx.ResponseSpec{
				Target:     domain.ErrGroupNotFound,
				Status:     http.StatusNotFound,
				Message:    "group not found",
				LogMessage: "group not found",
				LogLevel:   zapcore.WarnLevel,
			},
			httpx.ResponseSpec{
				Target:  domain.ErrForbidden,
				Status:  http.StatusForbidden,
				Message: "access forbidden",
			},
		) {
			return
		}

		response := mapGroupToUpdateResponse(group)

		c.JSON(http.StatusOK, response)
	}
}
