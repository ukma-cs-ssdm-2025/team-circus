package document

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/requests"
	"go.uber.org/zap"
)

type createDocumentService interface {
	Create(ctx context.Context, userUUID, groupUUID uuid.UUID, name, content string) (*domain.Document, error)
}

// NewCreateDocumentHandler creates a new document
// @Summary Create a new document
// @Description Create a new document with the provided group UUID, name and content
// @Tags documents
// @Accept json
// @Produce json
// @Param request body requests.CreateDocumentRequest true "Document creation request"
// @Success 201 {object} responses.CreateDocumentResponse "Document created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation failed"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents [post]
func NewCreateDocumentHandler(service createDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, exists := c.Get("user_role"); exists {
			if roleStr, ok := role.(string); ok && roleStr == domain.RoleViewer {
				c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
				return
			}
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

		var req requests.CreateDocumentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("create document handler: failed to bind request: %v", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("create document handler: validation failed: %v", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		document, err := service.Create(c.Request.Context(), userUUID, req.GroupUUID, req.Name, req.Content)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to create document",
				zap.Error(err),
				zap.String("group_uuid", req.GroupUUID.String()),
				zap.String("name", req.Name),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document"})
			return
		}
		if err != nil {
			logger.Error("failed to create document",
				zap.Error(err),
				zap.String("group_uuid", req.GroupUUID.String()),
				zap.String("name", req.Name),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document"})
			return
		}

		response := mapDocumentToCreateResponse(document)

		c.JSON(http.StatusCreated, response)
	}
}
