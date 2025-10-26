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
		return nil, domain.ErrDocumentNotFound
	}

	return document, nil
}

func (s *DocumentService) GetByUUIDForUser(ctx context.Context, documentUUID, userUUID uuid.UUID) (*domain.Document, error) {
	document, err := s.GetByUUID(ctx, documentUUID)
	if err != nil {
		return nil, err
	}

	isMember, err := s.groupRepo.IsMember(ctx, document.GroupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("document service: getByUUIDForUser: %w", err)
	}

	if !isMember {
		return nil, domain.ErrForbidden
	}

	return document, nil
}

func (s *DocumentService) GetAll(ctx context.Context) ([]*domain.Document, error) {
	documents, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("document service: getAll: %w", err)
	}

	return documents, nil
}

func (s *DocumentService) GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Document, error) {
	documents, err := s.repo.GetAllForUser(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("document service: getAllForUser: %w", err)
	}

	return documents, nil
}
