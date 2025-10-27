package group

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

type stubGroupMemberRepo struct {
	getByUUIDFn            func(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error)
	getMemberFn            func(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error)
	listMembersFn          func(ctx context.Context, groupUUID uuid.UUID) ([]*domain.GroupMember, error)
	addMemberFn            func(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) (*domain.GroupMember, error)
	updateMemberRoleFn     func(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) error
	removeMemberFn         func(ctx context.Context, groupUUID, userUUID uuid.UUID) error
	countMembersWithRoleFn func(ctx context.Context, groupUUID uuid.UUID, role string) (int, error)
}

func (s *stubGroupMemberRepo) GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error) {
	if s.getByUUIDFn != nil {
		return s.getByUUIDFn(ctx, groupUUID)
	}
	panic("unexpected call to GetByUUID")
}

func (s *stubGroupMemberRepo) GetMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error) {
	if s.getMemberFn != nil {
		return s.getMemberFn(ctx, groupUUID, userUUID)
	}
	panic("unexpected call to GetMember")
}

func (s *stubGroupMemberRepo) ListMembers(ctx context.Context, groupUUID uuid.UUID) ([]*domain.GroupMember, error) {
	if s.listMembersFn != nil {
		return s.listMembersFn(ctx, groupUUID)
	}
	panic("unexpected call to ListMembers")
}

func (s *stubGroupMemberRepo) AddMember(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) (*domain.GroupMember, error) {
	if s.addMemberFn != nil {
		return s.addMemberFn(ctx, groupUUID, userUUID, role)
	}
	panic("unexpected call to AddMember")
}

func (s *stubGroupMemberRepo) UpdateMemberRole(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) error {
	if s.updateMemberRoleFn != nil {
		return s.updateMemberRoleFn(ctx, groupUUID, userUUID, role)
	}
	panic("unexpected call to UpdateMemberRole")
}

func (s *stubGroupMemberRepo) RemoveMember(ctx context.Context, groupUUID, userUUID uuid.UUID) error {
	if s.removeMemberFn != nil {
		return s.removeMemberFn(ctx, groupUUID, userUUID)
	}
	panic("unexpected call to RemoveMember")
}

func (s *stubGroupMemberRepo) CountMembersWithRole(ctx context.Context, groupUUID uuid.UUID, role string) (int, error) {
	if s.countMembersWithRoleFn != nil {
		return s.countMembersWithRoleFn(ctx, groupUUID, role)
	}
	panic("unexpected call to CountMembersWithRole")
}

type stubUserRepo struct {
	getByUUIDFn func(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
}

func (s *stubUserRepo) GetByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if s.getByUUIDFn != nil {
		return s.getByUUIDFn(ctx, id)
	}
	panic("unexpected call to GetByUUID")
}

func TestGroupMemberService_ListMembers_GroupNotFound(t *testing.T) {
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return nil, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	members, err := service.ListMembers(context.Background(), uuid.New(), uuid.New())

	require.ErrorIs(t, err, domain.ErrGroupNotFound)
	require.Nil(t, members)
}

func TestGroupMemberService_ListMembers_ForbiddenForNonMember(t *testing.T) {
	groupID := uuid.New()
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			require.Equal(t, groupID, id)
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			require.Equal(t, groupID, gid)
			return nil, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	members, err := service.ListMembers(context.Background(), uuid.New(), groupID)

	require.ErrorIs(t, err, domain.ErrForbidden)
	require.Nil(t, members)
}

func TestGroupMemberService_ListMembers_Success(t *testing.T) {
	groupID := uuid.New()
	requesterID := uuid.New()
	expected := []*domain.GroupMember{
		{
			GroupUUID: groupID,
			UserUUID:  requesterID,
			Role:      domain.GroupRoleAuthor,
		},
	}

	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			return &domain.GroupMember{
				GroupUUID: gid,
				UserUUID:  uid,
				Role:      domain.GroupRoleAuthor,
			}, nil
		},
		listMembersFn: func(ctx context.Context, id uuid.UUID) ([]*domain.GroupMember, error) {
			require.Equal(t, groupID, id)
			return expected, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	members, err := service.ListMembers(context.Background(), requesterID, groupID)

	require.NoError(t, err)
	require.Equal(t, expected, members)
}

func TestGroupMemberService_AddMember_GroupNotFound(t *testing.T) {
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return nil, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	member, err := service.AddMember(context.Background(), uuid.New(), uuid.New(), uuid.New(), domain.GroupRoleReviewer)

	require.ErrorIs(t, err, domain.ErrGroupNotFound)
	require.Nil(t, member)
}
