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
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Group, error)
}

type getAllGroupsService interface {
	GetAll(ctx context.Context) ([]*domain.Group, error)
}

// NewGetGroupHandler retrieves a group by UUID
// @Summary Get a group by UUID
// @Description Retrieve a specific group by its UUID
// @Tags groups
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Success 200 {object} responses.GetGroupResponse "Group retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups/{uuid} [get]
func NewGetGroupHandler(service getGroupService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("get group handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		group, err := service.GetByUUID(c, parsedUUID)
		if errors.Is(err, domain.ErrGroupNotFound) {
			logger.Warn("group not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
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
// @Description Retrieve a list of all groups
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllGroupsResponse "Groups retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /groups [get]
func NewGetAllGroupsHandler(service getAllGroupsService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		groups, err := service.GetAll(c)
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
