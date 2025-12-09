package websocket

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Message type constants for Yjs protocol (y-protocols spec)
const (
	YjsSyncStep1         = 0
	YjsSyncStep2         = 1
	YjsUpdate            = 2
	MessageTypeAwareness = 101
)

// ClientConnection represents active WebSocket connection
type ClientConnection struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	UserName   string
	DocumentID uuid.UUID
	Conn       *websocket.Conn
	Send       chan []byte
	Done       chan struct{}
	LastSeen   time.Time
	// AwarenessID is the client id assigned by y-protocols/awareness
	AwarenessID uint32
	// CanEdit controls whether the client is allowed to apply document updates.
	CanEdit bool
}

// DocumentHub manages all connections for single document
type DocumentHub struct {
	DocumentID  uuid.UUID
	Clients     map[*ClientConnection]bool
	YjsDoc      []byte
	Broadcast   chan []byte
	Register    chan *ClientConnection
	Unregister  chan *ClientConnection
	Done        chan struct{}
	LastUpdated time.Time
	Version     int
}

// PersistenceRecord for database storage
type PersistenceRecord struct {
	DocumentID     uuid.UUID
	YjsSnapshot    []byte
	Version        int
	LastModifiedBy uuid.UUID
	LastModified   time.Time
}
