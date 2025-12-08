package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *MemberRepository) CreateMember(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) (*domain.Member, error) {
	const query = `
		INSERT INTO user_groups (group_uuid, user_uuid, role)
		VALUES ($1, $2, $3)
		RETURNING group_uuid, user_uuid, role, created_at`

	var member domain.Member
	err := r.db.QueryRowContext(ctx, query, groupUUID, userUUID, role).Scan(
		&member.GroupUUID,
		&member.UserUUID,
		&member.Role,
		&member.CreatedAt,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("member repository: create member: %w", err))
	}

	return &member, nil
}
