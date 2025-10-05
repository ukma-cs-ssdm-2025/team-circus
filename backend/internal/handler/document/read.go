package document

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/responses"
)

type getDocumentService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Document, error)
}

type getAllDocumentsService interface {
	GetAll(ctx context.Context) ([]*domain.Document, error)
}

type getDocumentsByGroupService interface {
	GetByGroupUUID(ctx context.Context, groupUUID uuid.UUID) ([]*domain.Document, error)
}

// NewGetDocumentHandler retrieves a document by UUID
// @Summary Get a document by UUID
// @Description Retrieve a specific document by its UUID
// @Tags documents
// @Accept json
// @Produce json
// @Param uuid path string true "Document UUID"
// @Success 200 {object} responses.GetDocumentResponse "Document retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 404 {object} map[string]interface{} "Document not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/documents/{uuid} [get]
func NewGetDocumentHandler(service getDocumentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("get document handler: failed to parse uuid: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		document, err := service.GetByUUID(c, parsedUUID)
		if errors.Is(err, domain.ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get document"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get document"})
			return
		}

		response := mapDocumentToGetResponse(document)

		c.JSON(http.StatusOK, response)
	}
}

// NewGetAllDocumentsHandler retrieves all documents
// @Summary Get all documents
// @Description Retrieve a list of all documents
// @Tags documents
// @Accept json
// @Produce json
// @Success 200 {object} responses.GetAllDocumentsResponse "Documents retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/documents [get]
func NewGetAllDocumentsHandler(service getAllDocumentsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		documents, err := service.GetAll(c)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
			return
		}

		response := responses.GetAllDocumentsResponse{
			Documents: mapDocumentsToGetAllResponse(documents),
		}

		c.JSON(http.StatusOK, response)
	}
}

// NewGetDocumentsByGroupHandler retrieves documents by group UUID
// @Summary Get documents by group UUID
// @Description Retrieve all documents for a specific group
// @Tags documents
// @Accept json
// @Produce json
// @Param uuid path string true "Group UUID"
// @Success 200 {object} responses.GetDocumentsByGroupResponse "Documents retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid UUID format"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/groups/{uuid}/documents [get]
func NewGetDocumentsByGroupHandler(service getDocumentsByGroupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupUUIDParam := c.Param("uuid")
		parsedGroupUUID, err := uuid.Parse(groupUUIDParam)
		if err != nil {
			err = fmt.Errorf("get documents by group handler: failed to parse group uuid: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group uuid format"})
			return
		}

		documents, err := service.GetByGroupUUID(c, parsedGroupUUID)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
			return
		}

		response := responses.GetDocumentsByGroupResponse{
			Documents: mapDocumentsToGetAllResponse(documents),
		}

		c.JSON(http.StatusOK, response)
	}
}
