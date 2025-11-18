package member

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/member/requests"
	"go.uber.org/zap"
)

type updateMemberService interface {
	UpdateMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID, role string) (*domain.Member, error)
}

func NewUpdateMemberHandler(service updateMemberService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupUUIDParam := c.Param("uuid")
		groupUUID, err := uuid.Parse(groupUUIDParam)
		if err != nil {
			err = fmt.Errorf("update member handler: failed to parse group uuid: %w", err)
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
			err = fmt.Errorf("update member handler: failed to parse user uuid: %w", err)
			logger.Error("failed to parse user uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user uuid format"})
			return
		}

		var req requests.UpdateMemberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("update member handler: failed to bind request: %w", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("update member handler: validation failed: %w", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		member, err := service.UpdateMemberByUser(c.Request.Context(),
			userUUID, groupUUID, memberUUID, req.Role)
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change the last author"})
			return
		}
		if err != nil {
			logger.Error("failed to update member role", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update member"})
			return
		}

		response := mapMemberToResponse(member)

		c.JSON(http.StatusOK, response)
	}
}
