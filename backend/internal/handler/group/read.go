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

type getGroupService interface {
	GetByUUIDForUser(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.Group, error)
}

type getAllGroupsService interface {
	GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error)
}

// NewGetGroupHandler retrieves a group by UUID
// @Summary Get a group by UUID
// @Description Retrieve a specific group by its UUID if the requesting user is a member
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Success 200 {object} responses.GetGroupResponse "Group retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid} [get]
func NewGetGroupHandler(service getGroupService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		groupUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"get group handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
			return
		}

		group, err := service.GetByUUIDForUser(c.Request.Context(), groupUUID, userUUID)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to get group",
				LogMessage: "failed to get group",
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

		response := mapGroupToGetResponse(group)

		c.JSON(http.StatusOK, response)
	}
}

// NewGetAllGroupsHandler retrieves all groups
// @Summary Get all groups
// @Description Retrieve a list of all groups the requesting user belongs to
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllGroupsResponse "Groups retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups [get]
func NewGetAllGroupsHandler(service getAllGroupsService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		groups, err := service.GetAllForUser(c.Request.Context(), userUUID)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to get groups",
				LogMessage: "failed to get groups",
				LogLevel:   zapcore.ErrorLevel,
			},
			httpx.RequestContextFields(c),
		) {
			return
		}

		response := responses.GetAllGroupsResponse{
			Groups: mapGroupsToGetAllResponse(groups),
		}

		c.JSON(http.StatusOK, response)
	}
}
