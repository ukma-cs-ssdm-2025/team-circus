package document

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *DocumentRepository) Delete(ctx context.Context, userUUID, documentUUID uuid.UUID) error {
	query := `
        DELETE FROM documents d
        USING user_groups ug
        WHERE d.uuid = $2
          AND ug.group_uuid = d.group_uuid
          AND ug.user_uuid = $1
          AND ug.role <> $3`

	result, err := r.db.ExecContext(ctx, query, userUUID, documentUUID, domain.GroupRoleReviewer)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("document repository: delete exec: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("document repository: delete rows affected: %w", err))
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
