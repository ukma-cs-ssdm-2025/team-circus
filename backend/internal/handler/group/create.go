package group

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/requests"
)

type createGroupService interface {
	Create(ctx context.Context, name string) (*domain.Group, error)
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
func NewCreateGroupHandler(service createGroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.CreateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("create group handler: failed to bind request: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("create group handler: validation failed: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		group, err := service.Create(c, req.Name)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
			return
		}

		response := mapGroupToCreateResponse(group)

		c.JSON(http.StatusCreated, response)
	}
}
