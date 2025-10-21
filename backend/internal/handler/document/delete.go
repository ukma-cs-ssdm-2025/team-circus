package document

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/responses"
	"go.uber.org/zap"
)

type deleteDocumentService interface {
	Delete(ctx context.Context, uuid uuid.UUID) error
}

// NewDeleteDocumentHandler deletes a document by UUID
// @Summary Delete a document by UUID
// @Description Delete a specific document by its UUID
// @Tags documents
// @Accept json
// @Produce json
// @Param uuid path string true "Document UUID"
// @Success 200 {object} responses.DeleteDocumentResponse "Document deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 404 {object} map[string]interface{} "Document not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents/{uuid} [delete]
func NewDeleteDocumentHandler(service deleteDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("delete document handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		err = service.Delete(c, parsedUUID)
		if errors.Is(err, domain.ErrDocumentNotFound) {
			logger.Warn("document not found", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to delete document", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete document"})
			return
		}
		if err != nil {
			logger.Error("failed to delete document", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete document"})
			return
		}

		response := responses.DeleteDocumentResponse{
			Message: "Document deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
