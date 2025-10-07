package user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

func (r *UserRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM users WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
