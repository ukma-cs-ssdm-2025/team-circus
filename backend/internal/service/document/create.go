package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Create(ctx context.Context, userUUID, groupUUID uuid.UUID, name, content string) (*domain.Document, error) {
	if err := s.ensureCanEditDocuments(ctx, groupUUID, userUUID); err != nil {
		return nil, err
	}

	document, err := s.repo.Create(ctx, groupUUID, name, content)
	if err != nil {
		return nil, fmt.Errorf("document service: create: %w", err)
	}

	return document, nil
}
