package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Document, error) {
	document, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("document service: getByUUID: %w", err)
	}

	if document == nil {
		return nil, ErrDocumentNotFound
	}

	return document, nil
}

func (s *DocumentService) GetByGroupUUID(ctx context.Context, groupUUID uuid.UUID) ([]*domain.Document, error) {
	documents, err := s.repo.GetByGroupUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("document service: getByGroupUUID: %w", err)
	}

	return documents, nil
}

func (s *DocumentService) GetAll(ctx context.Context) ([]*domain.Document, error) {
	documents, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("document service: getAll: %w", err)
	}

	return documents, nil
}
