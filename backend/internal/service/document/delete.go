package document

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Delete(ctx context.Context, docUUID, userUUID uuid.UUID) error {
	doc, err := s.GetByUUID(ctx, docUUID)
	if err != nil {
		return err
	}

	member, err := s.memberRepo.GetMember(ctx, doc.GroupUUID, userUUID)
	if err != nil {
		return fmt.Errorf("document service: delete: %w", err)
	}
	if member == nil {
		return domain.ErrForbidden
	}
	if member.Role == domain.RoleViewer {
		return domain.ErrForbidden
	}

	err = s.repo.Delete(ctx, docUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrDocumentNotFound
		}
		return fmt.Errorf("document service: delete: %w", err)
	}

	return nil
}
