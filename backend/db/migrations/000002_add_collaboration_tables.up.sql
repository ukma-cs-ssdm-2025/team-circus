-- Snapshot table: stores complete Y.Doc state
CREATE TABLE IF NOT EXISTS document_snapshots (
    document_id UUID PRIMARY KEY REFERENCES documents(uuid) ON DELETE CASCADE,
    yjs_snapshot BYTEA NOT NULL,
    version INTEGER DEFAULT 1,
    modified_by UUID REFERENCES users(uuid),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_document_snapshots_created_at ON document_snapshots(created_at);

-- Updates table: stores incremental changes (for audit trail & recovery)
CREATE TABLE IF NOT EXISTS document_updates (
    id BIGSERIAL PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(uuid) ON DELETE CASCADE,
    yjs_update BYTEA NOT NULL,
    user_id UUID NOT NULL REFERENCES users(uuid),
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_document_updates_document_id ON document_updates(document_id);
CREATE INDEX idx_document_updates_created_at ON document_updates(created_at);
CREATE INDEX idx_document_updates_document_id_version ON document_updates(document_id, version);

-- Presence table: track who's currently editing
CREATE TABLE IF NOT EXISTS document_presence (
    id BIGSERIAL PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(uuid) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(uuid),
    cursor_position INTEGER,
    last_seen TIMESTAMP DEFAULT NOW(),
    UNIQUE (document_id, user_id)
);

CREATE INDEX idx_document_presence_document_id ON document_presence(document_id);
