package member

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/member/responses"
	"go.uber.org/zap"
)

type getAllMembersService interface {
	GetAllMembersForUser(ctx context.Context, userUUID, groupUUID uuid.UUID) ([]*domain.Member, error)
}

func NewGetAllMembersHandler(service getAllMembersService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupUUIDParam := c.Param("uuid")
		groupUUID, err := uuid.Parse(groupUUIDParam)
		if err != nil {
			err = fmt.Errorf("get all members handler: failed to parse uuid: %w", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

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

		members, err := service.GetAllMembersForUser(c.Request.Context(), userUUID, groupUUID)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if err != nil {
			logger.Error("failed to list group members", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members"})
			return
		}

		response := responses.GetAllMembersResponse{
			Members: mapMembersToResponse(members),
		}

		c.JSON(http.StatusOK, response)
	}
}
