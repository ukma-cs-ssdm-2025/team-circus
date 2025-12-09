package document

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/responses"
	documentservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/document"
	"go.uber.org/zap"
)

type shareDocumentService interface {
	GenerateShareLink(ctx context.Context, docUUID, userUUID uuid.UUID, expirationDays int) (string, time.Time, error)
}

type publicDocumentService interface {
	ValidateShareLink(ctx context.Context, docParam, sigParam, expParam string) (*domain.Document, error)
}

// @Summary Generate a shareable link for a document
// @Description Create a time-limited, read-only share link for the specified document
// @Tags documents
// @Accept json
// @Produce json
// @Param uuid path string true "Document UUID"
// @Param request body requests.ShareDocumentRequest false "Share link options"
// @Success 200 {object} responses.ShareDocumentResponse "Share link generated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Authentication required"
// @Failure 403 {object} map[string]interface{} "Access forbidden"
// @Failure 404 {object} map[string]interface{} "Document not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents/{uuid}/share [post]
func NewShareDocumentHandler(service shareDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, exists := c.Get("user_role"); exists {
			if roleStr, ok := role.(string); ok && roleStr == domain.RoleViewer {
				c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
				return
			}
		}

		uuidParam := c.Param("uuid")
		docUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			err = fmt.Errorf("share document handler: failed to parse uuid: %v", err)
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

		var req requests.ShareDocumentRequest
		if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
			logger.Error("failed to bind share request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		shareURL, expiresAt, err := service.GenerateShareLink(
			c.Request.Context(),
			docUUID,
			userUUID,
			req.ExpirationDays,
		)
		switch {
		case errors.Is(err, domain.ErrDocumentNotFound):
			logger.Warn("document not found for share", zap.String("uuid", uuidParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		case errors.Is(err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		case errors.Is(err, documentservice.ErrInvalidExpiration):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expiration range"})
			return
		case err != nil:
			logger.Error("failed to generate share link", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate share link"})
			return
		}

		response := responses.ShareDocumentResponse{
			DocumentUUID: docUUID,
			URL:          shareURL,
			ExpiresAt:    expiresAt,
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Access a shared document
// @Description Retrieve a document using a public share link. No authentication required.
// @Tags documents
// @Accept json
// @Produce json
// @Param doc query string true "Document UUID from share link"
// @Param sig query string true "HMAC signature"
// @Param exp query string true "Expiration timestamp (unix seconds)"
// @Success 200 {object} responses.GetDocumentResponse "Document retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Missing or invalid parameters"
// @Failure 404 {object} map[string]interface{} "Invalid share link"
// @Failure 410 {object} map[string]interface{} "Share link expired"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /documents/public [get]
func NewGetPublicDocumentHandler(service publicDocumentService, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		docParam := c.Query("doc")
		sigParam := c.Query("sig")
		expParam := c.Query("exp")

		if docParam == "" || sigParam == "" || expParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
			return
		}

		document, err := service.ValidateShareLink(c.Request.Context(), docParam, sigParam, expParam)
		switch {
		case errors.Is(err, domain.ErrShareLinkExpired):
			c.JSON(http.StatusGone, gin.H{"error": "share link expired"})
			return
		case errors.Is(err, domain.ErrShareLinkInvalid):
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid share link"})
			return
		case errors.Is(err, domain.ErrDocumentNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		case err != nil:
			logger.Error("failed to validate share link", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load document"})
			return
		}

		c.Set("user_role", domain.RoleViewer)
		response := mapDocumentToGetResponse(document)

		c.JSON(http.StatusOK, response)
	}
}
