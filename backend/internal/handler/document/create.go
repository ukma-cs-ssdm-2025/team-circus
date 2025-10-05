package document

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/requests"
)

type createDocumentService interface {
	Create(ctx context.Context, groupUUID string, name, content string) (*domain.Document, error)
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
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/documents [post]
func NewCreateDocumentHandler(service createDocumentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req requests.CreateDocumentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("create document handler: failed to bind request: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := req.Validate(); err != nil {
			err = fmt.Errorf("create document handler: validation failed: %v", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		document, err := service.Create(c, req.GroupUUID.String(), req.Name, req.Content)
		if errors.Is(err, domain.ErrInternal) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document"})
			return
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create document"})
			return
		}

		response := mapDocumentToCreateResponse(document)

		c.JSON(http.StatusCreated, response)
	}
}
