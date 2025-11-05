package group

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/responses"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/httpx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type deleteGroupService interface {
	Delete(ctx context.Context, userUUID, groupUUID uuid.UUID) error
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
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		groupUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"delete group handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
			return
		}

		err := service.Delete(c.Request.Context(), userUUID, groupUUID)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to delete group",
				LogMessage: "failed to delete group",
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

		response := responses.DeleteGroupResponse{
			Message: "Group deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
