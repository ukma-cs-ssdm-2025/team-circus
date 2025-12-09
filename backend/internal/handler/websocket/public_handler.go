package websocket

import (
	"context"
	"errors"
	"hash/crc32"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"go.uber.org/zap"
)

type publicDocumentService interface {
	ValidateShareLink(ctx context.Context, docParam, sigParam, expParam string) (*domain.Document, error)
}

// NewPublicWebSocketHandler upgrades WebSocket connections for shared links (read-only guests).
func NewPublicWebSocketHandler(
	shareService publicDocumentService,
	hubManager *HubManager,
	logger *zap.Logger,
	allowedOrigins []string,
) gin.HandlerFunc {
	upgrader := createUpgrader(allowedOrigins)

	return func(c *gin.Context) {
		docParam := c.Param("uuid")
		sigParam := c.Query("sig")
		expParam := c.Query("exp")

		if docParam == "" || sigParam == "" || expParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
			return
		}

		document, err := shareService.ValidateShareLink(c.Request.Context(), docParam, sigParam, expParam)
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
			logger.Error("failed to validate share link for websocket", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open websocket"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("failed to upgrade public websocket connection", zap.Error(err))
			return
		}

		awarenessID, err := generateAwarenessID()
		if err != nil {
			logger.Warn("failed to generate awareness id for guest, falling back to timestamp", zap.Error(err))
			awarenessID = crc32.ChecksumIEEE([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
		}

		client := &ClientConnection{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			UserName:    "guest",
			DocumentID:  document.UUID,
			Conn:        conn,
			Send:        make(chan []byte, initialSendBufferSize),
			Done:        make(chan struct{}),
			LastSeen:    time.Now(),
			AwarenessID: awarenessID,
			CanEdit:     false,
		}

		hub := hubManager.GetOrCreateHub(document.UUID)

		select {
		case hub.Register <- client:
		case <-c.Request.Context().Done():
			if err := conn.Close(); err != nil {
				logger.Warn("failed to close websocket connection on context cancellation (public)",
					zap.Error(err),
					zap.String("document_id", document.UUID.String()),
					zap.String("path", c.Request.URL.Path),
				)
			}
			return
		}

		go writePump(client, logger)
		readPump(c.Request.Context(), hubManager, hub, client, logger)
	}
}
