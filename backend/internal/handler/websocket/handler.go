package websocket

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 65536
)

func generateAwarenessID() (uint32, error) {
	var buf [4]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

// createUpgrader creates a websocket upgrader with secure origin checking.
func createUpgrader(allowedOrigins []string) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")

			// Deny empty or missing Origin header by default
			if origin == "" {
				return false
			}

			// Never allow wildcard in production
			for _, allowed := range allowedOrigins {
				if allowed == "*" {
					// Only allow "*" if explicitly configured (not recommended for production)
					return true
				}

				// Exact match
				if origin == allowed {
					return true
				}

				// Support wildcard subdomain matching (e.g., "*.example.com")
				if strings.HasPrefix(allowed, "*.") {
					domain := allowed[2:] // Remove "*."
					if strings.HasSuffix(origin, domain) {
						// Ensure it's actually a subdomain match, not just a suffix match
						if origin == "http://"+domain || origin == "https://"+domain ||
							strings.HasSuffix(origin, "."+domain) {
							return true
						}
					}
				}
			}

			return false
		},
	}
}

type documentAccessService interface {
	GetByUUIDForUser(ctx context.Context, documentUUID, userUUID uuid.UUID) (*domain.Document, error)
	GetMemberRole(ctx context.Context, documentUUID, userUUID uuid.UUID) (string, error)
}

// NewWebSocketHandler upgrades HTTP connections to collaborative WebSocket sessions.
func NewWebSocketHandler(
	documentService documentAccessService,
	hubManager *HubManager,
	logger *zap.Logger,
	allowedOrigins []string,
) gin.HandlerFunc {
	upgrader := createUpgrader(allowedOrigins)

	return func(c *gin.Context) {
		docParam := c.Param("uuid")
		documentID, err := uuid.Parse(docParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
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

		_, err = documentService.GetByUUIDForUser(c.Request.Context(), documentID, userUUID)
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if errors.Is(err, domain.ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if err != nil {
			logger.Error("failed to get document for websocket", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open websocket"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("failed to upgrade connection", zap.Error(err))
			return
		}

		awarenessID, err := generateAwarenessID()
		if err != nil {
			logger.Warn("failed to generate awareness id, falling back to timestamp", zap.Error(err))
			awarenessID = crc32.ChecksumIEEE([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
		}

		role, err := documentService.GetMemberRole(c.Request.Context(), documentID, userUUID)
		if err != nil {
			if errors.Is(err, domain.ErrForbidden) {
				c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
				return
			}
			if !errors.Is(err, domain.ErrDocumentNotFound) {
				logger.Error("failed to resolve member role for websocket client", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open websocket"})
				return
			}
		}
		canEdit := role != domain.RoleViewer

		client := &ClientConnection{
			ID:          uuid.New(),
			UserID:      userUUID,
			UserName:    userUUID.String(),
			DocumentID:  documentID,
			Conn:        conn,
			Send:        make(chan []byte, initialSendBufferSize),
			Done:        make(chan struct{}),
			LastSeen:    time.Now(),
			AwarenessID: awarenessID,
			CanEdit:     canEdit,
		}

		hub := hubManager.GetOrCreateHub(documentID)
		select {
		case hub.Register <- client:
		case <-c.Request.Context().Done():
			if err := conn.Close(); err != nil {
				logger.Warn("failed to close websocket connection on context cancellation",
					zap.Error(err),
					zap.String("document_id", documentID.String()),
					zap.String("user_id", userUUID.String()),
					zap.String("path", c.Request.URL.Path),
				)
			}
			return
		}

		go writePump(client, logger)
		readPump(c.Request.Context(), hubManager, hub, client, logger)
	}
}

func writePump(client *ClientConnection, logger *zap.Logger) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			logger.Debug("failed to write close message in writePump defer", zap.Error(err))
		}
		if err := client.Conn.Close(); err != nil {
			logger.Debug("failed to close websocket in writePump defer", zap.Error(err))
		}
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				if err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					logger.Debug("failed to write close message when channel closed", zap.Error(err))
				}
				return
			}

			if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logger.Warn("failed to set write deadline for binary message",
					zap.Error(err),
					zap.String("client_id", client.ID.String()),
					zap.String("user_id", client.UserID.String()),
					zap.String("document_id", client.DocumentID.String()),
				)
				return
			}
			if err := client.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				logger.Warn("failed to write websocket message", zap.Error(err))
				return
			}
		case <-ticker.C:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logger.Warn("failed to set write deadline for ping message",
					zap.Error(err),
					zap.String("client_id", client.ID.String()),
					zap.String("user_id", client.UserID.String()),
					zap.String("document_id", client.DocumentID.String()),
				)
				return
			}
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-client.Done:
			return
		}
	}
}

func readPump(
	requestCtx context.Context,
	hubManager *HubManager,
	hub *DocumentHub,
	client *ClientConnection,
	logger *zap.Logger,
) {
	defer func() {
		select {
		case hub.Unregister <- client:
		default:
		}
		close(client.Done)
		if err := client.Conn.Close(); err != nil {
			logger.Debug("failed to close websocket in readPump defer", zap.Error(err))
		}
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logger.Warn("failed to set initial read deadline",
			zap.Error(err),
			zap.String("client_id", client.ID.String()),
			zap.String("user_id", client.UserID.String()),
			zap.String("document_id", client.DocumentID.String()),
		)
		return
	}
	client.Conn.SetPongHandler(func(string) error {
		if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			logger.Warn("failed to set read deadline in pong handler",
				zap.Error(err),
				zap.String("client_id", client.ID.String()),
				zap.String("user_id", client.UserID.String()),
				zap.String("document_id", client.DocumentID.String()),
			)
			return err
		}
		client.LastSeen = time.Now()
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Warn("unexpected websocket close", zap.Error(err))
			}
			break
		}

		handleClientMessage(requestCtx, hubManager, hub, client, logger, message)

		select {
		case <-requestCtx.Done():
			return
		default:
		}
	}
}

func handleAwareness(hub *DocumentHub, client *ClientConnection, message []byte) {
	client.LastSeen = time.Now()
	select {
	case hub.Broadcast <- message:
	default:
	}
}

func handleClientMessage(
	requestCtx context.Context,
	hubManager *HubManager,
	hub *DocumentHub,
	client *ClientConnection,
	logger *zap.Logger,
	message []byte,
) {
	client.LastSeen = time.Now()
	if len(message) == 0 {
		return
	}

	msgType := message[0]
	switch msgType {
	case MessageTypeAwareness:
		handleAwareness(hub, client, message)
	default:
		if msgType == YjsUpdate && !client.CanEdit {
			logger.Debug(
				"blocking update from read-only client",
				zap.String("document_id", hub.DocumentID.String()),
				zap.String("client_id", client.ID.String()),
			)
			return
		}

		select {
		case hub.Broadcast <- message:
		default:
			logger.Warn(
				"dropping broadcast message; channel full",
				zap.String("document_id", hub.DocumentID.String()),
			)
		}

		if msgType == YjsUpdate && hubManager.persistence != nil {
			versionToStore := hub.Version
			if versionToStore < 1 {
				versionToStore = defaultHubVersion
			}
			if err := hubManager.persistence.SaveUpdate(
				requestCtx,
				hub.DocumentID,
				client.UserID,
				message,
				versionToStore+1,
			); err != nil {
				logger.Warn(
					"failed to persist update",
					zap.Error(err),
					zap.String("document_id", hub.DocumentID.String()),
				)
			}
		}
	}
}
