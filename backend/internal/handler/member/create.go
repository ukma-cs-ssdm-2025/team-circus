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

type createMemberService interface {
	CreateMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID,
		role string) (*domain.Member, error)
}

func NewCreateMemberHandler(service createMemberService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupUUIDParam := c.Param("uuid")
		groupUUID, err := uuid.Parse(groupUUIDParam)
		if err != nil {
			err = fmt.Errorf("create member handler: failed to parse group uuid: %w", err)
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

		var req requests.CreateMemberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("create member handler: failed to bind request: %w", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("create member handler: validation failed: %w", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		member, err := service.CreateMemberByUser(c.Request.Context(),
			userUUID, groupUUID, req.UserUUID, req.Role)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrGroupNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, domain.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "member already exists"})
			return
		}
		if errors.Is(err, domain.ErrOnlyAuthor) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "there must be only one author"})
			return
		}
		if err != nil {
			logger.Error("failed to create member", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create member"})
			return
		}

		response := mapMemberToResponse(member)

		c.JSON(http.StatusCreated, response)
	}
}
