package document

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/responses"
	"go.uber.org/zap"
)

const (
	wsWriteWait     = 10 * time.Second
	wsPongWait      = 60 * time.Second
	wsPingPeriod    = (wsPongWait * 9) / 10
	wsMaxMessageLen = int64(1 << 20) // 1MB

	wsMessageTypeInit   = "init"
	wsMessageTypeUpdate = "update"
	wsMessageTypeError  = "error"
)

type documentRealtimeService interface {
	GetByUUIDForUser(ctx context.Context, documentUUID, userUUID uuid.UUID) (*domain.Document, error)
	Update(ctx context.Context, docUUID, userUUID uuid.UUID, name, content string) (*domain.Document, error)
}

type documentUpdateMessage struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type documentSocketMessage struct {
	Type     string                            `json:"type"`
	Document *responses.UpdateDocumentResponse `json:"document,omitempty"`
	Error    string                            `json:"error,omitempty"`
}

type DocumentHub struct {
	mu     sync.RWMutex
	rooms  map[uuid.UUID]map[*documentClient]struct{}
	logger *zap.Logger
}

func NewDocumentHub(logger *zap.Logger) *DocumentHub {
	return &DocumentHub{
		rooms:  make(map[uuid.UUID]map[*documentClient]struct{}),
		logger: logger,
	}
}

func (h *DocumentHub) register(client *documentClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[client.documentUUID]; !ok {
		h.rooms[client.documentUUID] = make(map[*documentClient]struct{})
	}
	h.rooms[client.documentUUID][client] = struct{}{}
}

func (h *DocumentHub) unregister(client *documentClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[client.documentUUID]
	if !ok {
		return
	}

	if _, exists := room[client]; exists {
		delete(room, client)
		close(client.send)
	}

	if len(room) == 0 {
		delete(h.rooms, client.documentUUID)
	}
}

func (h *DocumentHub) BroadcastDocument(document *domain.Document) {
	if document == nil {
		return
	}

	payload, err := encodeDocumentMessage(document, wsMessageTypeUpdate)
	if err != nil {
		h.logger.Error("failed to marshal websocket document message", zap.Error(err))
		return
	}

	h.mu.RLock()
	clients := make([]*documentClient, 0, len(h.rooms[document.UUID]))
	for client := range h.rooms[document.UUID] {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		client.enqueue(payload)
	}
}

type documentClient struct {
	conn         *websocket.Conn
	send         chan []byte
	hub          *DocumentHub
	documentUUID uuid.UUID
	userUUID     uuid.UUID
	logger       *zap.Logger
	ctx          context.Context
	cancel       context.CancelFunc
	closeOnce    sync.Once
}

func newDocumentClient(conn *websocket.Conn, hub *DocumentHub, documentUUID, userUUID uuid.UUID, logger *zap.Logger) *documentClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &documentClient{
		conn:         conn,
		send:         make(chan []byte, 8),
		hub:          hub,
		documentUUID: documentUUID,
		userUUID:     userUUID,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (c *documentClient) enqueue(payload []byte) {
	select {
	case c.send <- payload:
	default:
		c.logger.Warn("websocket send buffer full, closing client", zap.String("document_uuid", c.documentUUID.String()))
		c.close()
	}
}

func (c *documentClient) sendError(message string) {
	payload, err := json.Marshal(documentSocketMessage{
		Type:  wsMessageTypeError,
		Error: message,
	})
	if err != nil {
		c.logger.Error("failed to marshal websocket error message", zap.Error(err))
		return
	}
	c.enqueue(payload)
}

func (c *documentClient) close() {
	c.closeOnce.Do(func() {
		c.cancel()
		c.hub.unregister(c)
		_ = c.conn.Close()
	})
}

func (c *documentClient) readPump(service documentRealtimeService) {
	defer c.close()

	c.conn.SetReadLimit(wsMaxMessageLen)
	_ = c.conn.SetReadDeadline(time.Now().Add(wsPongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(wsPongWait))
		return nil
	})

	for {
		var msg documentUpdateMessage
		if err := c.conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Warn("websocket read error", zap.Error(err), zap.String("document_uuid", c.documentUUID.String()))
			}
			break
		}

		if msg.Type != wsMessageTypeUpdate {
			continue
		}

		req := requests.UpdateDocumentRequest{
			Name:    msg.Name,
			Content: msg.Content,
		}
		if err := req.Validate(); err != nil {
			c.sendError("validation failed: " + err.Error())
			continue
		}

		updatedDoc, err := service.Update(c.ctx, c.documentUUID, c.userUUID, req.Name, req.Content)
		if errors.Is(err, domain.ErrDocumentNotFound) {
			c.sendError("document not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			c.sendError("access forbidden")
			return
		}
		if errors.Is(err, domain.ErrInternal) {
			c.sendError("failed to update document")
			continue
		}
		if err != nil {
			c.sendError("failed to update document")
			continue
		}

		c.hub.BroadcastDocument(updatedDoc)
	}
}

func (c *documentClient) writePump() {
	ticker := time.NewTicker(wsPingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = c.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func encodeDocumentMessage(document *domain.Document, messageType string) ([]byte, error) {
	if document == nil {
		return nil, errors.New("document is nil")
	}

	docResponse := mapDocumentToUpdateResponse(document)

	return json.Marshal(documentSocketMessage{
		Type:     messageType,
		Document: &docResponse,
	})
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return true
	}
	if len(allowedOrigins) == 0 {
		return true
	}
	for _, allowed := range allowedOrigins {
		if allowed == "*" || strings.EqualFold(strings.TrimRight(allowed, "/"), strings.TrimRight(origin, "/")) {
			return true
		}
	}
	return false
}

func NewDocumentWebsocketHandler(service documentRealtimeService, hub *DocumentHub, logger *zap.Logger, allowedOrigins []string) gin.HandlerFunc {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isOriginAllowed(r.Header.Get("Origin"), allowedOrigins)
		},
	}

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
		docUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			logger.Error("websocket: failed to parse uuid", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		document, err := service.GetByUUIDForUser(context.Background(), docUUID, userUUID)
		if errors.Is(err, domain.ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			return
		}
		if err != nil {
			logger.Error("websocket: failed to load document", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load document"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("websocket upgrade failed", zap.Error(err))
			return
		}

		client := newDocumentClient(conn, hub, docUUID, userUUID, logger)
		hub.register(client)

		initialPayload, err := encodeDocumentMessage(document, wsMessageTypeInit)
		if err != nil {
			logger.Error("failed to marshal initial websocket message", zap.Error(err))
			client.close()
			return
		}
		client.enqueue(initialPayload)

		go client.writePump()
		go client.readPump(service)
	}
}

type broadcastingUpdateService struct {
	service updateDocumentService
	hub     *DocumentHub
}

func NewBroadcastingUpdateService(service updateDocumentService, hub *DocumentHub) updateDocumentService {
	if hub == nil {
		return service
	}
	return &broadcastingUpdateService{
		service: service,
		hub:     hub,
	}
}

func (s *broadcastingUpdateService) Update(ctx context.Context, docUUID, userUUID uuid.UUID, name, content string) (*domain.Document, error) {
	document, err := s.service.Update(ctx, docUUID, userUUID, name, content)
	if err != nil || document == nil {
		return document, err
	}

	s.hub.BroadcastDocument(document)
	return document, nil
}
