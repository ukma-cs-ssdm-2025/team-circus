package document

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (s *DocumentService) Delete(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.Delete(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrDocumentNotFound
		}
		return fmt.Errorf("document service: delete: %w", err)
	}

	return nil
}
