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

type deleteMemberService interface {
	DeleteMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID) error
}

func NewDeleteMemberHandler(service deleteMemberService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupUUIDParam := c.Param("uuid")
		groupUUID, err := uuid.Parse(groupUUIDParam)
		if err != nil {
			err = fmt.Errorf("delete member handler: failed to parse group uuid: %w", err)
			logger.Error("failed to parse group uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group uuid format"})
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

		memberUUIDParam := c.Param("user_uuid")
		memberUUID, err := uuid.Parse(memberUUIDParam)
		if err != nil {
			err = fmt.Errorf("delete member handler: failed to parse user uuid: %w", err)
			logger.Error("failed to parse user uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user uuid format"})
			return
		}

		err = service.DeleteMemberByUser(c.Request.Context(), userUUID, groupUUID, memberUUID)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		if errors.Is(err, domain.ErrOnlyAuthor) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot remove the last author"})
			return
		}
		if err != nil {
			logger.Error("failed to remove group member", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove member"})
			return
		}

		response := responses.DeleteMemberResponse{
			Message: "Member deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
