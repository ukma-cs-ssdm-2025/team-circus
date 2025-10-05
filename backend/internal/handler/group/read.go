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

type getGroupService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Group, error)
}

type getAllGroupsService interface {
	GetAll(ctx context.Context) ([]*domain.Group, error)
}

func NewGetGroupHandler(service getGroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("get group handler: failed to parse uuid: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		group, err := service.GetByUUID(c, parsedUUID)
		if errors.Is(err, domain.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group"})
			return
		}

		response := mapGroupToGetResponse(group)

		c.JSON(http.StatusOK, response)
	}
}

func NewGetAllGroupsHandler(service getAllGroupsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		groups, err := service.GetAll(c)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get groups"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get groups"})
			return
		}

		response := responses.GetAllGroupsResponse{
			Groups: mapGroupsToGetAllResponse(groups),
		}

		c.JSON(http.StatusOK, response)
	}
}
