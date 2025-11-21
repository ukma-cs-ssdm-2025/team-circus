package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Update(ctx context.Context, docUUID, userUUID uuid.UUID, name, content string) (*domain.Document, error) {
	doc, err := s.GetByUUID(ctx, docUUID)
	if err != nil {
		return nil, err
	}

	member, err := s.memberRepo.GetMember(ctx, doc.GroupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("document service: update: %w", err)
	}
	if member == nil {
		return nil, domain.ErrForbidden
	}
	if member.Role == domain.RoleViewer {
		return nil, domain.ErrForbidden
	}

	updatedDoc, err := s.repo.Update(ctx, docUUID, name, content)
	if err != nil {
		return nil, fmt.Errorf("document service: update: %w", err)
	}
	if updatedDoc == nil {
		return nil, domain.ErrDocumentNotFound
	}

	return updatedDoc, nil
}
