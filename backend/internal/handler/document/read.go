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

type getDocumentService interface {
	GetByUUIDForUser(ctx context.Context, documentUUID, userUUID uuid.UUID) (*domain.Document, error)
}

type getAllDocumentsService interface {
	GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Document, error)
}

// NewGetDocumentHandler retrieves a document by UUID
// @Summary Get a document by UUID
// @Description Retrieve a specific document by its UUID if the requesting user is a member of the owning group
// @Tags documents
// @Accept json
// @Produce json
// @Param uuid path string true "Document UUID"
// @Success 200 {object} responses.GetDocumentResponse "Document retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Document not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents/{uuid} [get]
func NewGetDocumentHandler(service getDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		documentUUID, ok := httpx.ParseUUIDParam(
			c,
			logger,
			"uuid",
			"get document handler: failed to parse uuid",
			httpx.RequestContextFields(c)...,
		)
		if !ok {
			return
		}

		document, err := service.GetByUUIDForUser(c.Request.Context(), documentUUID, userUUID)
		if handleDocumentOperationError(c, logger, err, documentUUID, "get") {
			return
		}

		response := mapDocumentToGetResponse(document)

		c.JSON(http.StatusOK, response)
	}
}

// NewGetAllDocumentsHandler retrieves all documents
// @Summary Get all documents
// @Description Retrieve a list of all documents belonging to groups the requesting user is a member of
// @Tags documents
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllDocumentsResponse "Documents retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents [get]
func NewGetAllDocumentsHandler(service getAllDocumentsService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
			return
		}

		documents, err := service.GetAllForUser(c.Request.Context(), userUUID)
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to get documents",
				LogMessage: "failed to get documents",
				LogLevel:   zapcore.ErrorLevel,
			},
			httpx.RequestContextFields(c),
		) {
			return
		}

		response := responses.GetAllDocumentsResponse{
			Documents: mapDocumentsToGetAllResponse(documents),
		}

		c.JSON(http.StatusOK, response)
	}
}
