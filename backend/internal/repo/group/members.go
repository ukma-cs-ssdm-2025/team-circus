package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *GroupRepository) GetMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error) {
	const query = `
		SELECT ug.group_uuid, ug.user_uuid, ug.role, ug.created_at, u.login, u.email
		FROM user_groups ug
		INNER JOIN users u ON u.uuid = ug.user_uuid
		WHERE ug.group_uuid = $1 AND ug.user_uuid = $2`

	var member domain.GroupMember
	err := r.db.QueryRowContext(ctx, query, groupUUID, userUUID).Scan(
		&member.GroupUUID,
		&member.UserUUID,
		&member.Role,
		&member.CreatedAt,
		&member.UserLogin,
		&member.UserEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: get member: %w", err))
	}

	return &member, nil
}

func (r *GroupRepository) ListMembers(ctx context.Context, groupUUID uuid.UUID) ([]*domain.GroupMember, error) {
	const query = `
		SELECT ug.group_uuid, ug.user_uuid, ug.role, ug.created_at, u.login, u.email
		FROM user_groups ug
		INNER JOIN users u ON u.uuid = ug.user_uuid
		WHERE ug.group_uuid = $1
		ORDER BY u.login`

	rows, err := r.db.QueryContext(ctx, query, groupUUID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: list members query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var members []*domain.GroupMember
	for rows.Next() {
		var member domain.GroupMember
		err := rows.Scan(
			&member.GroupUUID,
			&member.UserUUID,
			&member.Role,
			&member.CreatedAt,
			&member.UserLogin,
			&member.UserEmail,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: list members scan: %w", err))
		}
		members = append(members, &member)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: list members rows err: %w", err))
	}

	return members, nil
}

func (r *GroupRepository) AddMember(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) (*domain.GroupMember, error) {
	const query = `
		INSERT INTO user_groups (user_uuid, group_uuid, role)
		VALUES ($1, $2, $3)
		RETURNING group_uuid, user_uuid, role, created_at`

	var member domain.GroupMember
	if err := r.db.QueryRowContext(ctx, query, userUUID, groupUUID, role).Scan(
		&member.GroupUUID,
		&member.UserUUID,
		&member.Role,
		&member.CreatedAt,
	); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: add member: %w", err))
	}

	return &member, nil
}

func (r *GroupRepository) UpdateMemberRole(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) error {
	const query = `
		UPDATE user_groups
		SET role = $3
		WHERE group_uuid = $1 AND user_uuid = $2`

	result, err := r.db.ExecContext(ctx, query, groupUUID, userUUID, role)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member role exec: %w", err))
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member role rows: %w", err))
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *GroupRepository) RemoveMember(ctx context.Context, groupUUID, userUUID uuid.UUID) error {
	const query = `
		DELETE FROM user_groups
		WHERE group_uuid = $1 AND user_uuid = $2`

	result, err := r.db.ExecContext(ctx, query, groupUUID, userUUID)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: remove member exec: %w", err))
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: remove member rows: %w", err))
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *GroupRepository) CountMembersWithRole(ctx context.Context, groupUUID uuid.UUID, role string) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM user_groups
		WHERE group_uuid = $1 AND role = $2`

	var count int
	if err := r.db.QueryRowContext(ctx, query, groupUUID, role).Scan(&count); err != nil {
		return 0, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: count members with role: %w", err))
	}

	return count, nil
}
