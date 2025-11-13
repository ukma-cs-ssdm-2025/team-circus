package document

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/httpx"
	"go.uber.org/zap"
)

type updateDocumentService interface {
	Update(ctx context.Context, userUUID, documentUUID uuid.UUID, name, content string) (*domain.Document, error)
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
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		documentUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"update document handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
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

		document, err := service.Update(c.Request.Context(), userUUID, documentUUID, req.Name, req.Content)
		if handleDocumentOperationError(c, logger, err, documentUUID, "update") {
			return
		}

		response := mapDocumentToUpdateResponse(document)

		c.JSON(http.StatusOK, response)
	}
}
