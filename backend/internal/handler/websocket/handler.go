package websocket

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type documentAccessService interface {
	GetByUUIDForUser(ctx context.Context, documentUUID, userUUID uuid.UUID) (*domain.Document, error)
}

// NewWebSocketHandler upgrades HTTP connections to collaborative WebSocket sessions.
func NewWebSocketHandler(
	documentService documentAccessService,
	hubManager *HubManager,
	logger *zap.Logger,
) gin.HandlerFunc {
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

		client := &ClientConnection{
			ID:          uuid.New(),
			UserID:      userUUID,
			UserName:    userUUID.String(),
			DocumentID:  documentID,
			Conn:        conn,
			Send:        make(chan []byte, initialSendBufferSize),
			Done:        make(chan struct{}),
			LastSeen:    time.Now(),
			AwarenessID: rand.New(rand.NewSource(time.Now().UnixNano())).Uint32(),
		}

		hub := hubManager.GetOrCreateHub(documentID)
		select {
		case hub.Register <- client:
		case <-c.Request.Context().Done():
			_ = conn.Close()
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
		_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		_ = client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				logger.Warn("failed to write websocket message", zap.Error(err))
				return
			}
		case <-ticker.C:
			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
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
		_ = client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	_ = client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		_ = client.Conn.SetReadDeadline(time.Now().Add(pongWait))
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

		client.LastSeen = time.Now()
		if len(message) == 0 {
			continue
		}

		msgType := message[0]
		switch msgType {
		case MessageTypeAwareness:
			handleAwareness(hub, client, message)
		default:
			select {
			case hub.Broadcast <- message:
			default:
				logger.Warn("dropping broadcast message; channel full", zap.String("document_id", hub.DocumentID.String()))
			}

			if msgType == YjsUpdate && hubManager.persistence != nil {
				versionToStore := hub.Version
				if versionToStore < 1 {
					versionToStore = defaultHubVersion
				}
				err = hubManager.persistence.SaveUpdate(context.Background(), hub.DocumentID, client.UserID, message, versionToStore+1)
				if err != nil {
					logger.Warn("failed to persist update", zap.Error(err), zap.String("document_id", hub.DocumentID.String()))
				}
			}
		}

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
