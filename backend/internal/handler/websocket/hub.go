package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo"
	"go.uber.org/zap"
)

const (
	defaultHubVersion     = 1
	persistenceInterval   = 60 * time.Second
	registerBufferSize    = 32
	unregisterBufferSize  = 32
	broadcastBufferSize   = 128
	initialSendBufferSize = 128
)

// HubManager coordinates hubs per document and handles persistence.
type HubManager struct {
	hubs        map[uuid.UUID]*DocumentHub
	mu          sync.RWMutex
	logger      *zap.Logger
	persistence *repo.DocumentPersistence
}

// NewHubManager constructs a manager for collaboration hubs.
func NewHubManager(logger *zap.Logger, persistence *repo.DocumentPersistence) *HubManager {
	return &HubManager{
		hubs:        make(map[uuid.UUID]*DocumentHub),
		logger:      logger,
		persistence: persistence,
	}
}

// GetOrCreateHub returns an existing hub for a document or initializes a new one, loading persisted state if available.
func (m *HubManager) GetOrCreateHub(documentID uuid.UUID) *DocumentHub {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, ok := m.hubs[documentID]; ok {
		return hub
	}

	hub := &DocumentHub{
		DocumentID:  documentID,
		Clients:     make(map[*ClientConnection]bool),
		Broadcast:   make(chan []byte, broadcastBufferSize),
		Register:    make(chan *ClientConnection, registerBufferSize),
		Unregister:  make(chan *ClientConnection, unregisterBufferSize),
		Done:        make(chan struct{}),
		Version:     defaultHubVersion,
		LastUpdated: time.Now(),
	}

	if m.persistence != nil {
		record, err := m.persistence.LoadSnapshot(context.Background(), documentID)
		if err != nil {
			m.logger.Warn("failed to load snapshot for document", zap.String("document_id", documentID.String()), zap.Error(err))
		} else if record != nil {
			hub.YjsDoc = record.YjsSnapshot
			hub.Version = record.Version
			hub.LastUpdated = record.LastModified
		}
	}

	m.hubs[documentID] = hub
	go m.run(hub)
	return hub
}

// CloseHub terminates a hub and removes it from manager.
func (m *HubManager) CloseHub(documentID uuid.UUID) {
	m.mu.Lock()
	hub, ok := m.hubs[documentID]
	if ok {
		delete(m.hubs, documentID)
	}
	m.mu.Unlock()

	if ok {
		close(hub.Done)
	}
}

func (m *HubManager) run(hub *DocumentHub) {
	ticker := time.NewTicker(persistenceInterval)
	defer ticker.Stop()

	for {
		select {
		case client := <-hub.Register:
			m.registerClient(hub, client)
		case client := <-hub.Unregister:
			m.unregisterClient(hub, client)
		case message := <-hub.Broadcast:
			m.broadcastMessage(hub, message)
		case <-ticker.C:
			m.persistHubState(hub)
		case <-hub.Done:
			m.cleanupHub(hub)
			return
		}
	}
}

func (m *HubManager) registerClient(hub *DocumentHub, client *ClientConnection) {
	hub.Clients[client] = true
	client.LastSeen = time.Now()

	if len(hub.YjsDoc) > 0 {
		select {
		case client.Send <- hub.YjsDoc:
		default:
			m.logger.Warn(
				"dropping initial state, client send buffer full",
				zap.String("document_id", hub.DocumentID.String()),
				zap.String("user_id", client.UserID.String()),
			)
		}
	}
}

func (m *HubManager) unregisterClient(hub *DocumentHub, client *ClientConnection) {
	if _, ok := hub.Clients[client]; ok {
		delete(hub.Clients, client)
		close(client.Send)
	}

	if len(hub.Clients) == 0 {
		m.persistHubState(hub)
		m.CloseHub(hub.DocumentID)
	}
}

func (m *HubManager) broadcastMessage(hub *DocumentHub, message []byte) {
	hub.LastUpdated = time.Now()

	if len(message) > 0 && message[0] == YjsUpdate {
		hub.Version++
		hub.YjsDoc = message
	}

	for client := range hub.Clients {
		select {
		case client.Send <- message:
		default:
			// Drop messages to slow clients to keep hub responsive.
			close(client.Send)
			delete(hub.Clients, client)
		}
	}
}

func (m *HubManager) persistHubState(hub *DocumentHub) {
	if m.persistence == nil {
		return
	}

	if len(hub.YjsDoc) == 0 {
		return
	}

	err := m.persistence.SaveSnapshot(context.Background(), hub.DocumentID, hub.YjsDoc, hub.Version, uuid.Nil)
	if err != nil {
		m.logger.Warn("failed to persist snapshot", zap.String("document_id", hub.DocumentID.String()), zap.Error(err))
	}
}

func (m *HubManager) cleanupHub(hub *DocumentHub) {
	for client := range hub.Clients {
		close(client.Send)
		delete(hub.Clients, client)
	}
}
