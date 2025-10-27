package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Update(ctx context.Context, userUUID, documentUUID uuid.UUID, name, content string) (*domain.Document, error) {
	current, err := s.repo.GetByUUID(ctx, documentUUID)
	if err != nil {
		return nil, fmt.Errorf("document service: update get document: %w", err)
	}

	if current == nil {
		return nil, domain.ErrDocumentNotFound
	}

	if err := s.ensureCanEditDocuments(ctx, current.GroupUUID, userUUID); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, documentUUID, name, content)
	if err != nil {
		return nil, fmt.Errorf("document service: update: %w", err)
	}

	if updated == nil {
		return nil, domain.ErrDocumentNotFound
	}

	return updated, nil
}
