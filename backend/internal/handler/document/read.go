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

		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("get document handler: failed to parse uuid: %v", err)
			logger.Error("failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		document, err := service.GetByUUIDForUser(c, parsedUUID, userUUID)
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
			logger.Error("failed to get document", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get document"})
			return
		}
		if err != nil {
			logger.Error("failed to get document", zap.Error(err), zap.String("uuid", uuidParam))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get document"})
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
		userUUIDValue, exists := c.Get("user_uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user context missing"})
			return
		}

		userUUID, ok := userUUIDValue.(uuid.UUID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
			return
		}

		documents, err := service.GetAllForUser(c, userUUID)
		if errors.Is(err, domain.ErrInternal) {
			logger.Error("failed to get documents", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
			return
		}
		if err != nil {
			logger.Error("failed to get documents", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
			return
		}

		response := responses.GetAllDocumentsResponse{
			Documents: mapDocumentsToGetAllResponse(documents),
		}

		c.JSON(http.StatusOK, response)
	}
}
