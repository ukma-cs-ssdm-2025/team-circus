package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

// SnapshotRecord mirrors the persisted snapshot for a collaborative document.
type SnapshotRecord struct {
	DocumentID     uuid.UUID
	YjsSnapshot    []byte
	Version        int
	LastModifiedBy uuid.UUID
	LastModified   time.Time
}

// UpdateRecord represents an incremental CRDT update row.
type UpdateRecord struct {
	ID         int64
	DocumentID uuid.UUID
	YjsUpdate  []byte
	UserID     uuid.UUID
	Version    int
	CreatedAt  time.Time
}

// PresenceRecord represents a user's current presence in a document.
type PresenceRecord struct {
	ID             int64
	DocumentID     uuid.UUID
	UserID         uuid.UUID
	CursorPosition *int
	LastSeen       time.Time
}

// DocumentPersistence stores and retrieves collaborative state.
type DocumentPersistence struct {
	db *sql.DB
}

// NewDocumentPersistence constructs a persistence helper.
func NewDocumentPersistence(db *sql.DB) *DocumentPersistence {
	return &DocumentPersistence{db: db}
}

// SaveSnapshot upserts the latest snapshot for a document.
func (p *DocumentPersistence) SaveSnapshot(ctx context.Context, documentID uuid.UUID, yjsSnapshot []byte, version int, modifiedBy uuid.UUID) error {
	query := `
		INSERT INTO document_snapshots (document_id, yjs_snapshot, version, modified_by, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (document_id) DO UPDATE
		SET yjs_snapshot = EXCLUDED.yjs_snapshot,
		    version = EXCLUDED.version,
		    modified_by = EXCLUDED.modified_by,
		    updated_at = NOW()
	`

	modifiedByUUID := uuid.NullUUID{}
	if modifiedBy != uuid.Nil {
		modifiedByUUID = uuid.NullUUID{UUID: modifiedBy, Valid: true}
	}

	_, err := p.db.ExecContext(ctx, query, documentID, yjsSnapshot, version, modifiedByUUID)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: saveSnapshot: %w", err))
	}

	return nil
}

// LoadSnapshot retrieves the last persisted snapshot for a document.
func (p *DocumentPersistence) LoadSnapshot(ctx context.Context, documentID uuid.UUID) (*SnapshotRecord, error) {
	query := `
		SELECT document_id, yjs_snapshot, version, modified_by, updated_at
		FROM document_snapshots
		WHERE document_id = $1
	`

	var record SnapshotRecord
	var modifiedBy uuid.NullUUID
	err := p.db.QueryRowContext(ctx, query, documentID).Scan(
		&record.DocumentID,
		&record.YjsSnapshot,
		&record.Version,
		&modifiedBy,
		&record.LastModified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: loadSnapshot: %w", err))
	}

	if modifiedBy.Valid {
		record.LastModifiedBy = modifiedBy.UUID
	}

	return &record, nil
}

// SaveUpdate stores a CRDT incremental update for auditing or recovery.
func (p *DocumentPersistence) SaveUpdate(ctx context.Context, documentID, userID uuid.UUID, yjsUpdate []byte, version int) error {
	query := `
		INSERT INTO document_updates (document_id, yjs_update, user_id, version, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := p.db.ExecContext(ctx, query, documentID, yjsUpdate, userID, version)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: saveUpdate: %w", err))
	}

	return nil
}

// GetUpdates fetches updates since a provided version (exclusive).
func (p *DocumentPersistence) GetUpdates(ctx context.Context, documentID uuid.UUID, fromVersion int) ([]UpdateRecord, error) {
	query := `
		SELECT id, document_id, yjs_update, user_id, version, created_at
		FROM document_updates
		WHERE document_id = $1 AND version > $2
		ORDER BY id ASC
	`

	rows, err := p.db.QueryContext(ctx, query, documentID, fromVersion)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: getUpdates: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var updates []UpdateRecord
	for rows.Next() {
		var update UpdateRecord
		err = rows.Scan(
			&update.ID,
			&update.DocumentID,
			&update.YjsUpdate,
			&update.UserID,
			&update.Version,
			&update.CreatedAt,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: getUpdates scan: %w", err))
		}
		updates = append(updates, update)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: getUpdates rows: %w", err))
	}

	return updates, nil
}

// UpsertPresence updates or inserts a user's presence in a document.
// Uses ON CONFLICT to ensure only one presence row per (document_id, user_id).
func (p *DocumentPersistence) UpsertPresence(ctx context.Context, documentID, userID uuid.UUID, cursorPosition *int) error {
	query := `
		INSERT INTO document_presence (document_id, user_id, cursor_position, last_seen)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (document_id, user_id) DO UPDATE
		SET cursor_position = EXCLUDED.cursor_position,
		    last_seen = NOW()
	`

	_, err := p.db.ExecContext(ctx, query, documentID, userID, cursorPosition)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: upsertPresence: %w", err))
	}

	return nil
}

// GetPresence retrieves all active presence records for a document.
func (p *DocumentPersistence) GetPresence(ctx context.Context, documentID uuid.UUID) ([]PresenceRecord, error) {
	query := `
		SELECT id, document_id, user_id, cursor_position, last_seen
		FROM document_presence
		WHERE document_id = $1
		ORDER BY last_seen DESC
	`

	rows, err := p.db.QueryContext(ctx, query, documentID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: getPresence: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var presences []PresenceRecord
	for rows.Next() {
		var presence PresenceRecord
		var cursorPos sql.NullInt32
		err = rows.Scan(
			&presence.ID,
			&presence.DocumentID,
			&presence.UserID,
			&cursorPos,
			&presence.LastSeen,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: getPresence scan: %w", err))
		}

		if cursorPos.Valid {
			pos := int(cursorPos.Int32)
			presence.CursorPosition = &pos
		}

		presences = append(presences, presence)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: getPresence rows: %w", err))
	}

	return presences, nil
}

// RemovePresence removes a user's presence from a document (e.g., when they disconnect).
func (p *DocumentPersistence) RemovePresence(ctx context.Context, documentID, userID uuid.UUID) error {
	query := `
		DELETE FROM document_presence
		WHERE document_id = $1 AND user_id = $2
	`

	_, err := p.db.ExecContext(ctx, query, documentID, userID)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("document persistence: removePresence: %w", err))
	}

	return nil
}
