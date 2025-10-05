package group

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/responses"
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
// @Router /api/v1/groups/{uuid} [delete]
func NewDeleteGroupHandler(service deleteGroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("delete group handler: failed to parse uuid: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		err = service.Delete(c, parsedUUID)
		if errors.Is(err, domain.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group"})
			return
		}

		response := responses.DeleteGroupResponse{
			Message: "Group deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
