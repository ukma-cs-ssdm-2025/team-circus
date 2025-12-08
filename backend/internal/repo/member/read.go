package member

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *MemberRepository) GetMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.Member, error) {
	const query = `
		SELECT group_uuid, user_uuid, role, created_at
		FROM user_groups
		WHERE group_uuid = $1 AND user_uuid = $2`

	var member domain.Member
	err := r.db.QueryRowContext(ctx, query, groupUUID, userUUID).Scan(
		&member.GroupUUID,
		&member.UserUUID,
		&member.Role,
		&member.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: get member: %w", err))
	}

	return &member, nil
}

func (r *MemberRepository) GetAllMembers(ctx context.Context, groupUUID uuid.UUID) ([]*domain.Member, error) {
	const query = `
		SELECT group_uuid, user_uuid, role, created_at
		FROM user_groups
		WHERE group_uuid = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, groupUUID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: get all members query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var members []*domain.Member
	for rows.Next() {
		var member domain.Member
		err = rows.Scan(
			&member.GroupUUID,
			&member.UserUUID,
			&member.Role,
			&member.CreatedAt,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: get all members scan: %w", err))
		}
		members = append(members, &member)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: get all members rows err: %w", err))
	}

	return members, nil
}
