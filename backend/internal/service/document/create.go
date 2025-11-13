package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Create(ctx context.Context, userUUID, groupUUID uuid.UUID, name, content string) (*domain.Document, error) {
	isMember, err := s.groupRepo.IsMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("document service: create: %w", err)
	}

	if !isMember {
		return nil, domain.ErrForbidden
	}

	document, err := s.repo.Create(ctx, groupUUID, name, content)
	if err != nil {
		return nil, fmt.Errorf("document service: create: %w", err)
	}

	return document, nil
}
