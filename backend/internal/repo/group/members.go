package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, domain.ErrAlreadyExists
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: add member: %w", err))
	}

	return &member, nil
}

func (r *GroupRepository) UpdateMemberRole(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member role begin tx: %w", err))
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() //nolint:errcheck
		}
	}()

	const lockMemberQuery = `
		SELECT role
		FROM user_groups
		WHERE group_uuid = $1 AND user_uuid = $2
		FOR UPDATE`

	var currentRole string
	if scanErr := tx.QueryRowContext(ctx, lockMemberQuery, groupUUID, userUUID).Scan(&currentRole); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		err = errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member role select: %w", scanErr))
		return err
	}

	if currentRole == domain.GroupRoleAuthor && role != domain.GroupRoleAuthor {
		authorCount, countErr := countAuthorsForUpdate(ctx, tx, groupUUID)
		if countErr != nil {
			err = countErr
			return err
		}
		if authorCount <= 1 {
			err = domain.ErrLastAuthor
			return err
		}
	}

	const updateQuery = `
		UPDATE user_groups
		SET role = $3
		WHERE group_uuid = $1 AND user_uuid = $2`

	result, execErr := tx.ExecContext(ctx, updateQuery, groupUUID, userUUID, role)
	if execErr != nil {
		err = errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member role exec: %w", execErr))
		return err
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		err = errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member role rows: %w", rowsErr))
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update member role commit: %w", commitErr))
	}

	return nil
}

func (r *GroupRepository) RemoveMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: remove member begin tx: %w", err))
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() //nolint:errcheck
		}
	}()

	const lockMemberQuery = `
		SELECT role
		FROM user_groups
		WHERE group_uuid = $1 AND user_uuid = $2
		FOR UPDATE`

	var currentRole string
	if scanErr := tx.QueryRowContext(ctx, lockMemberQuery, groupUUID, userUUID).Scan(&currentRole); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		err = errors.Join(domain.ErrInternal, fmt.Errorf("group repository: remove member select: %w", scanErr))
		return err
	}

	if currentRole == domain.GroupRoleAuthor {
		authorCount, countErr := countAuthorsForUpdate(ctx, tx, groupUUID)
		if countErr != nil {
			err = countErr
			return err
		}
		if authorCount <= 1 {
			err = domain.ErrLastAuthor
			return err
		}
	}

	const deleteQuery = `
		DELETE FROM user_groups
		WHERE group_uuid = $1 AND user_uuid = $2`

	result, execErr := tx.ExecContext(ctx, deleteQuery, groupUUID, userUUID)
	if execErr != nil {
		err = errors.Join(domain.ErrInternal, fmt.Errorf("group repository: remove member exec: %w", execErr))
		return err
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		err = errors.Join(domain.ErrInternal, fmt.Errorf("group repository: remove member rows: %w", rowsErr))
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: remove member commit: %w", commitErr))
	}

	return nil
}

func countAuthorsForUpdate(ctx context.Context, tx *sql.Tx, groupUUID uuid.UUID) (int, error) {
	const query = `
		SELECT user_uuid
		FROM user_groups
		WHERE group_uuid = $1 AND role = $2
		FOR UPDATE`

	rows, err := tx.QueryContext(ctx, query, groupUUID, domain.GroupRoleAuthor)
	if err != nil {
		return 0, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: count authors for update query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	count := 0
	for rows.Next() {
		count++
	}

	if err := rows.Err(); err != nil {
		return 0, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: count authors for update rows: %w", err))
	}

	return count, nil
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
