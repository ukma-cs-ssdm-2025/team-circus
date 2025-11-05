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
	"go.uber.org/zap/zapcore"
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
// @Failure 404 {object} map[string]interface{} "Group not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents [post]
func NewCreateDocumentHandler(service createDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID, ok := httpx.ResolveUserUUID(c)
		if !ok {
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
		if httpx.HandleError(
			c,
			logger,
			err,
			httpx.ResponseSpec{
				Status:     http.StatusInternalServerError,
				Message:    "failed to create document",
				LogMessage: "failed to create document",
				LogLevel:   zapcore.ErrorLevel,
			},
			httpx.RequestContextFields(c, zap.String("group_uuid", req.GroupUUID.String()), zap.String("document_name", req.Name)),
			httpx.ResponseSpec{
				Target:  domain.ErrForbidden,
				Status:  http.StatusForbidden,
				Message: "access forbidden",
			},
			httpx.ResponseSpec{
				Target:  domain.ErrGroupNotFound,
				Status:  http.StatusNotFound,
				Message: "group not found",
			},
		) {
			return
		}

		response := mapDocumentToCreateResponse(document)

		c.JSON(http.StatusCreated, response)
	}
}
