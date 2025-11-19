package member

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *MemberRepository) DeleteMember(ctx context.Context, groupUUID, userUUID uuid.UUID) error {
	const query = `
		DELETE FROM user_groups
		WHERE group_uuid = $1 AND user_uuid = $2`

	result, err := r.db.ExecContext(ctx, query, groupUUID, userUUID)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: delete member exec: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: delete member rows affected: %w", err))
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
