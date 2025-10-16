package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *GroupRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM groups WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: delete exec: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: delete rows affected: %w", err))
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
