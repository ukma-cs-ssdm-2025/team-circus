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

func TestGroupMemberService_AddMember_ForbiddenForNonAuthorRequester(t *testing.T) {
	groupID := uuid.New()
	requesterID := uuid.New()
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == requesterID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleReviewer}, nil
			}
			return nil, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	member, err := service.AddMember(context.Background(), requesterID, groupID, uuid.New(), domain.GroupRoleReviewer)

	require.ErrorIs(t, err, domain.ErrForbidden)
	require.Nil(t, member)
}

func TestGroupMemberService_UpdateMemberRole_ForbiddenForNonAuthorRequester(t *testing.T) {
	groupID := uuid.New()
	requesterID := uuid.New()
	memberID := uuid.New()
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == requesterID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleReviewer}, nil
			}
			return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleReviewer}, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	updated, err := service.UpdateMemberRole(context.Background(), requesterID, groupID, memberID, domain.GroupRoleReviewer)

	require.ErrorIs(t, err, domain.ErrForbidden)
	require.Nil(t, updated)
}

func TestGroupMemberService_UpdateMemberRole_ForbiddenWhenModifyingAnotherAuthor(t *testing.T) {
	groupID := uuid.New()
	requesterID := uuid.New()
	otherAuthorID := uuid.New()
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == requesterID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
			}
			return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	updated, err := service.UpdateMemberRole(context.Background(), requesterID, groupID, otherAuthorID, domain.GroupRoleCoAuthor)

	require.ErrorIs(t, err, domain.ErrForbidden)
	require.Nil(t, updated)
}

func TestGroupMemberService_UpdateMemberRole_LastAuthorDemotionBlocked(t *testing.T) {
	groupID := uuid.New()
	authorID := uuid.New()
	memberRole := domain.GroupRoleAuthor
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == authorID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: memberRole}, nil
			}
			return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
		},
		updateMemberRoleFn: func(ctx context.Context, gid, uid uuid.UUID, role string) error {
			return domain.ErrLastAuthor
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	updated, err := service.UpdateMemberRole(context.Background(), authorID, groupID, authorID, domain.GroupRoleCoAuthor)

	require.ErrorIs(t, err, domain.ErrLastAuthor)
	require.Nil(t, updated)
}

func TestGroupMemberService_UpdateMemberRole_SelfDemoteWithMultipleAuthors(t *testing.T) {
	groupID := uuid.New()
	authorID := uuid.New()
	currentRole := domain.GroupRoleAuthor
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == authorID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: currentRole}, nil
			}
			return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
		},
		updateMemberRoleFn: func(ctx context.Context, gid, uid uuid.UUID, role string) error {
			require.Equal(t, authorID, uid)
			currentRole = role
			return nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	updated, err := service.UpdateMemberRole(context.Background(), authorID, groupID, authorID, domain.GroupRoleCoAuthor)

	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, domain.GroupRoleCoAuthor, updated.Role)
}

func TestGroupMemberService_RemoveMember_ForbiddenForNonAuthorRequester(t *testing.T) {
	groupID := uuid.New()
	requesterID := uuid.New()
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == requesterID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleReviewer}, nil
			}
			return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleReviewer}, nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	err := service.RemoveMember(context.Background(), requesterID, groupID, uuid.New())

	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestGroupMemberService_RemoveMember_AllowRemovingAnotherAuthorWhenMultipleAuthors(t *testing.T) {
	groupID := uuid.New()
	requesterID := uuid.New()
	otherAuthorID := uuid.New()
	removeCalled := false
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == requesterID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
			}
			if uid == otherAuthorID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
			}
			return nil, nil
		},
		countMembersWithRoleFn: func(ctx context.Context, gid uuid.UUID, role string) (int, error) {
			require.Equal(t, domain.GroupRoleAuthor, role)
			return 2, nil
		},
		removeMemberFn: func(ctx context.Context, gid, uid uuid.UUID) error {
			removeCalled = true
			return nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	err := service.RemoveMember(context.Background(), requesterID, groupID, otherAuthorID)

	require.NoError(t, err)
	require.True(t, removeCalled)
}

func TestGroupMemberService_RemoveMember_LastAuthorBlocked(t *testing.T) {
	groupID := uuid.New()
	authorID := uuid.New()
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == authorID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
			}
			return nil, nil
		},
		countMembersWithRoleFn: func(ctx context.Context, gid uuid.UUID, role string) (int, error) {
			require.Equal(t, domain.GroupRoleAuthor, role)
			return 1, nil
		},
		removeMemberFn: func(ctx context.Context, gid, uid uuid.UUID) error {
			return domain.ErrLastAuthor
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	err := service.RemoveMember(context.Background(), authorID, groupID, authorID)

	require.ErrorIs(t, err, domain.ErrLastAuthor)
}

func TestGroupMemberService_RemoveMember_Success(t *testing.T) {
	groupID := uuid.New()
	authorID := uuid.New()
	memberID := uuid.New()
	removeCalled := false
	repo := &stubGroupMemberRepo{
		getByUUIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
			return &domain.Group{UUID: id}, nil
		},
		getMemberFn: func(ctx context.Context, gid, uid uuid.UUID) (*domain.GroupMember, error) {
			if uid == authorID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleAuthor}, nil
			}
			if uid == memberID {
				return &domain.GroupMember{GroupUUID: gid, UserUUID: uid, Role: domain.GroupRoleReviewer}, nil
			}
			return nil, nil
		},
		countMembersWithRoleFn: func(ctx context.Context, gid uuid.UUID, role string) (int, error) {
			require.Equal(t, domain.GroupRoleAuthor, role)
			return 2, nil
		},
		removeMemberFn: func(ctx context.Context, gid, uid uuid.UUID) error {
			require.Equal(t, memberID, uid)
			removeCalled = true
			return nil
		},
	}
	service := NewGroupMemberService(repo, &stubUserRepo{})

	err := service.RemoveMember(context.Background(), authorID, groupID, memberID)

	require.NoError(t, err)
	require.True(t, removeCalled)
}

func TestValidateAssignableRole(t *testing.T) {
	testcases := []struct {
		name        string
		role        string
		expectedErr error
	}{
		{"coauthor valid", domain.GroupRoleCoAuthor, nil},
		{"reviewer valid", domain.GroupRoleReviewer, nil},
		{"author invalid", domain.GroupRoleAuthor, domain.ErrInvalidRole},
		{"unknown invalid", "unknown", domain.ErrInvalidRole},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAssignableRole(tc.role)
			if tc.expectedErr == nil {
				require.NoError(t, err)
				return
			}
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
