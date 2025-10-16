package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Update(ctx context.Context, uuid uuid.UUID, name, content string) (*domain.Document, error) {
	document, err := s.repo.Update(ctx, uuid, name, content)
	if err != nil {
		return nil, fmt.Errorf("document service: update: %w", err)
	}

	if document == nil {
		return nil, domain.ErrDocumentNotFound
	}

	return document, nil
}
