package member

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *MemberRepository) UpdateMember(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) (*domain.Member, error) {
	const query = `
		UPDATE user_groups
		SET role = $3
		WHERE group_uuid = $1 AND user_uuid = $2
		RETURNING group_uuid, user_uuid, role, created_at`

	var member domain.Member
	err := r.db.QueryRowContext(ctx, query, groupUUID, userUUID, role).Scan(
		&member.GroupUUID,
		&member.UserUUID,
		&member.Role,
		&member.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member: %w", err))
	}

	return &member, nil
}
