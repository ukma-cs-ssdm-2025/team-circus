package group

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

type stubGroupRepository struct {
	createFn        func(ctx context.Context, name string, authorUUID uuid.UUID) (*domain.Group, error)
	getByUUIDFn     func(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error)
	getMemberFn     func(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error)
	getAllFn        func(ctx context.Context) ([]*domain.Group, error)
	getAllForUserFn func(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error)
	updateFn        func(ctx context.Context, uuid uuid.UUID, name string) (*domain.Group, error)
	deleteFn        func(ctx context.Context, uuid uuid.UUID) error
}

func (s *stubGroupRepository) Create(ctx context.Context, name string, authorUUID uuid.UUID) (*domain.Group, error) {
	if s.createFn != nil {
		return s.createFn(ctx, name, authorUUID)
	}
	panic("unexpected call to Create")
}

func (s *stubGroupRepository) GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error) {
	if s.getByUUIDFn != nil {
		return s.getByUUIDFn(ctx, groupUUID)
	}
	panic("unexpected call to GetByUUID")
}

func (s *stubGroupRepository) GetMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error) {
	if s.getMemberFn != nil {
		return s.getMemberFn(ctx, groupUUID, userUUID)
	}
	panic("unexpected call to GetMember")
}

func (s *stubGroupRepository) GetAll(ctx context.Context) ([]*domain.Group, error) {
	if s.getAllFn != nil {
		return s.getAllFn(ctx)
	}
	panic("unexpected call to GetAll")
}

func (s *stubGroupRepository) GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error) {
	if s.getAllForUserFn != nil {
		return s.getAllForUserFn(ctx, userUUID)
	}
	panic("unexpected call to GetAllForUser")
}

func (s *stubGroupRepository) Update(ctx context.Context, id uuid.UUID, name string) (*domain.Group, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, name)
	}
	panic("unexpected call to Update")
}

func (s *stubGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	panic("unexpected call to Delete")
}

func TestGroupService_GetByUUIDForUser_GroupNotFound(t *testing.T) {
	repo := &stubGroupRepository{
		getByUUIDFn: func(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error) {
			return nil, nil
		},
	}
	service := NewGroupService(repo)

	group, err := service.GetByUUIDForUser(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, domain.ErrGroupNotFound)
	require.Nil(t, group)
}

func TestGroupService_GetByUUIDForUser_ForbiddenWhenNotMember(t *testing.T) {
	groupUUID := uuid.New()
	repo := &stubGroupRepository{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			require.Equal(t, groupUUID, id)
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, id, user uuid.UUID) (*domain.GroupMember, error) {
			require.Equal(t, groupUUID, id)
			return nil, nil
		},
	}
	service := NewGroupService(repo)

	group, err := service.GetByUUIDForUser(context.Background(), groupUUID, uuid.New())

	require.ErrorIs(t, err, domain.ErrForbidden)
	require.Nil(t, group)
}

func TestGroupService_GetByUUIDForUser_ReturnsGroupWithRole(t *testing.T) {
	groupUUID := uuid.New()
	userUUID := uuid.New()

	repo := &stubGroupRepository{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, id, user uuid.UUID) (*domain.GroupMember, error) {
			require.Equal(t, groupUUID, id)
			require.Equal(t, userUUID, user)
			return &domain.GroupMember{
				GroupUUID: id,
				UserUUID:  user,
				Role:      domain.GroupRoleAuthor,
			}, nil
		},
	}
	service := NewGroupService(repo)

	group, err := service.GetByUUIDForUser(context.Background(), groupUUID, userUUID)

	require.NoError(t, err)
	require.NotNil(t, group)
	require.Equal(t, domain.GroupRoleAuthor, group.Role)
	require.Equal(t, userUUID, group.AuthorUUID)
}
