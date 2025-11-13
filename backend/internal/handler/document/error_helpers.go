package document

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/httpx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// handleDocumentOperationError maps service errors to HTTP responses shared across document mutations.
func handleDocumentOperationError(
	c *gin.Context,
	logger *zap.Logger,
	err error,
	documentUUID uuid.UUID,
	operation string,
) bool {
	if err == nil {
		return false
	}

	return httpx.HandleError(
		c,
		logger,
		err,
		httpx.ResponseSpec{
			Status:     http.StatusInternalServerError,
			Message:    fmt.Sprintf("failed to %s document", operation),
			LogMessage: fmt.Sprintf("failed to %s document", operation),
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
	)
}
