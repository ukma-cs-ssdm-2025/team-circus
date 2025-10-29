package document

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *DocumentService) Delete(ctx context.Context, userUUID, documentUUID uuid.UUID) error {
	current, err := s.repo.GetByUUID(ctx, documentUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrDocumentNotFound
		}
		return fmt.Errorf("document service: delete get document: %w", err)
	}

	if current == nil {
		return domain.ErrDocumentNotFound
	}

	if err := s.ensureCanEditDocuments(ctx, current.GroupUUID, userUUID); err != nil {
		return err
	}

	err = s.repo.Delete(ctx, userUUID, documentUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrDocumentNotFound
		}
		return fmt.Errorf("document service: delete: %w", err)
	}

	return nil
}
