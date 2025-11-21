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

type updateDocumentService interface {
	Update(ctx context.Context, docUUID, userUUID uuid.UUID, name, content string) (*domain.Document, error)
}

// NewUpdateDocumentHandler updates a document by UUID
// @Summary Update a document by UUID
// @Description Update a specific document's name and content by its UUID
// @Tags documents
// @Accept json
// @Produce json
// @Param uuid path string true "Document UUID"
// @Param request body requests.UpdateDocumentRequest true "Document update request"
// @Success 200 {object} responses.UpdateDocumentResponse "Document updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format or validation failed"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Document not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents/{uuid} [put]
func NewUpdateDocumentHandler(service updateDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		docUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("update document handler: failed to parse uuid: %v", err)
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

		var req requests.UpdateDocumentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("update document handler: failed to bind request: %v", err)
			logger.Error("failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("update document handler: validation failed: %v", err)
			logger.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		document, err := service.Update(c.Request.Context(), docUUID, userUUID, req.Name, req.Content)
		if errors.Is(err, domain.ErrDocumentNotFound) {
			logger.Warn("document not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to update document",
				zap.Error(err),
				zap.String("uuid", uuidParam),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update document"})
			return
		}
		if err != nil {
			logger.Error("failed to update document",
				zap.Error(err),
				zap.String("uuid", uuidParam),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update document"})
			return
		}

		response := mapDocumentToUpdateResponse(document)

		c.JSON(http.StatusOK, response)
	}
}
