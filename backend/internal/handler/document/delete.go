package document

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/responses"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/httpx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type deleteDocumentService interface {
	Delete(ctx context.Context, userUUID, documentUUID uuid.UUID) error
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
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Document not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents/{uuid} [delete]
func NewDeleteDocumentHandler(service deleteDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		documentUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"delete document handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
			return
		}

		err := service.Delete(c.Request.Context(), userUUID, documentUUID)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to delete document",
				LogMessage: "failed to delete document",
				LogLevel:   zapcore.ErrorLevel,
			},
			httpx.RequestContextFields(c, zap.String("document_uuid", documentUUID.String())),
			httpx.ResponseSpec{
				Target:     domain.ErrDocumentNotFound,
				Status:     http.StatusNotFound,
				Message:    "document not found",
				LogMessage: "document not found",
				LogLevel:   zapcore.WarnLevel,
			},
			httpx.ResponseSpec{
				Target:  domain.ErrForbidden,
				Status:  http.StatusForbidden,
				Message: "access forbidden",
			},
		) {
			return
		}

		response := responses.DeleteDocumentResponse{
			Message: "Document deleted successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
