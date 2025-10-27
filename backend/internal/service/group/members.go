package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	userrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
)

type groupRepository interface {
	GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error)
	GetMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error)
	ListMembers(ctx context.Context, groupUUID uuid.UUID) ([]*domain.GroupMember, error)
	AddMember(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) (*domain.GroupMember, error)
	UpdateMemberRole(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, groupUUID, userUUID uuid.UUID) error
	CountMembersWithRole(ctx context.Context, groupUUID uuid.UUID, role string) (int, error)
}

type userRepository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
}

var (
	_ groupRepository = (*grouprepo.GroupRepository)(nil)
	_ userRepository  = (*userrepo.UserRepository)(nil)
)

type GroupMemberService struct {
	groupRepo groupRepository
	userRepo  userRepository
}

func NewGroupMemberService(groupRepo groupRepository, userRepo userRepository) *GroupMemberService {
	return &GroupMemberService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

func (s *GroupMemberService) ensureGroupMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error) {
	group, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: ensure group: %w", err)
	}
	if group == nil {
		return nil, domain.ErrGroupNotFound
	}

	member, err := s.groupRepo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: ensure membership: %w", err)
	}
	if member == nil {
		return nil, domain.ErrForbidden
	}

	return member, nil
}

func (s *GroupMemberService) ListMembers(ctx context.Context, requesterUUID, groupUUID uuid.UUID) ([]*domain.GroupMember, error) {
	if _, err := s.ensureGroupMember(ctx, groupUUID, requesterUUID); err != nil {
		return nil, fmt.Errorf("group member service: list members ensure requester: %w", err)
	}

	members, err := s.groupRepo.ListMembers(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: list members: %w", err)
	}

	return members, nil
}

func (s *GroupMemberService) AddMember(
	ctx context.Context,
	requesterUUID uuid.UUID,
	groupUUID uuid.UUID,
	memberUUID uuid.UUID,
	role string,
) (*domain.GroupMember, error) {
	if err := validateAssignableRole(role); err != nil {
		return nil, err
	}

	requesterMember, err := s.ensureGroupMember(ctx, groupUUID, requesterUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: add member requester: %w", err)
	}
	if requesterMember.Role != domain.GroupRoleAuthor {
		return nil, domain.ErrForbidden
	}

	if requesterUUID == memberUUID {
		return nil, domain.ErrAlreadyExists
	}

	user, err := s.userRepo.GetByUUID(ctx, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: add member get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	existingMember, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: add member get existing: %w", err)
	}
	if existingMember != nil {
		return nil, domain.ErrAlreadyExists
	}

	if _, err := s.groupRepo.AddMember(ctx, groupUUID, memberUUID, role); err != nil {
		return nil, fmt.Errorf("group member service: add member insert: %w", err)
	}

	createdMember, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: add member fetch created: %w", err)
	}

	return createdMember, nil
}

func (s *GroupMemberService) UpdateMemberRole(
	ctx context.Context,
	requesterUUID uuid.UUID,
	groupUUID uuid.UUID,
	memberUUID uuid.UUID,
	role string,
) (*domain.GroupMember, error) {
	if err := validateAssignableRole(role); err != nil {
		return nil, err
	}

	requesterMember, err := s.ensureGroupMember(ctx, groupUUID, requesterUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: update member requester: %w", err)
	}
	if requesterMember.Role != domain.GroupRoleAuthor {
		return nil, domain.ErrForbidden
	}

	member, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: update member get member: %w", err)
	}
	if member == nil {
		return nil, domain.ErrUserNotFound
	}

	if member.Role == domain.GroupRoleAuthor && memberUUID != requesterUUID {
		return nil, domain.ErrForbidden
	}

	if err := s.groupRepo.UpdateMemberRole(ctx, groupUUID, memberUUID, role); err != nil {
		if errors.Is(err, domain.ErrLastAuthor) {
			return nil, domain.ErrLastAuthor
		}
		return nil, fmt.Errorf("group member service: update member role: %w", err)
	}

	updatedMember, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: update member fetch updated: %w", err)
	}

	return updatedMember, nil
}

func (s *GroupMemberService) RemoveMember(ctx context.Context, requesterUUID, groupUUID, memberUUID uuid.UUID) error {
	requesterMember, err := s.ensureGroupMember(ctx, groupUUID, requesterUUID)
	if err != nil {
		return fmt.Errorf("group member service: remove member requester: %w", err)
	}
	if requesterMember.Role != domain.GroupRoleAuthor {
		return domain.ErrForbidden
	}

	member, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return fmt.Errorf("group member service: remove member get member: %w", err)
	}
	if member == nil {
		return domain.ErrUserNotFound
	}

	if err := s.groupRepo.RemoveMember(ctx, groupUUID, memberUUID); err != nil {
		if errors.Is(err, domain.ErrLastAuthor) {
			return domain.ErrLastAuthor
		}
		return fmt.Errorf("group member service: remove member delete: %w", err)
	}

	return nil
}

func validateAssignableRole(role string) error {
	switch role {
	case domain.GroupRoleCoAuthor, domain.GroupRoleReviewer:
		return nil
	case domain.GroupRoleAuthor:
		return domain.ErrInvalidRole
	default:
		return domain.ErrInvalidRole
	}
}
