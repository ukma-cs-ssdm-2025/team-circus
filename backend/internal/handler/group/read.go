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
			err = fmt.Errorf("get group handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		group, err := service.GetByUUIDForUser(c, parsedUUID, userUUID)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group"})
			return
		}
		if err != nil {
			logger.Error("failed to get group", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group"})
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
		userUUIDValue, exists := c.Get("user_uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user context missing"})
			return
		}

		userUUID, ok := userUUIDValue.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
			return
		}

		groups, err := service.GetAllForUser(c, userUUID)
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to get groups", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get groups"})
			return
		}
		if err != nil {
			logger.Error("failed to get groups", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get groups"})
			return
		}

		response := responses.GetAllGroupsResponse{
			Groups: mapGroupsToGetAllResponse(groups),
		}

		c.JSON(http.StatusOK, response)
	}
}
