package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Create(ctx context.Context, groupUUID uuid.UUID, name, content string) (*domain.Document, error) {
	document, err := s.repo.Create(ctx, groupUUID, name, content)
	if err != nil {
		return nil, fmt.Errorf("document service: create: %w", err)
	}

	return document, nil
}
