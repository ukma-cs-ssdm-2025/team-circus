package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	userrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
)

type GroupMemberService struct {
	groupRepo *grouprepo.GroupRepository
	userRepo  *userrepo.UserRepository
}

func NewGroupMemberService(groupRepo *grouprepo.GroupRepository, userRepo *userrepo.UserRepository) *GroupMemberService {
	return &GroupMemberService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

func (s *GroupMemberService) ListMembers(ctx context.Context, requesterUUID, groupUUID uuid.UUID) ([]*domain.GroupMember, error) {
	member, err := s.groupRepo.GetMember(ctx, groupUUID, requesterUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: list members get requester: %w", err)
	}
	if member == nil {
		return nil, domain.ErrForbidden
	}

	members, err := s.groupRepo.ListMembers(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: list members: %w", err)
	}

	return members, nil
}

func (s *GroupMemberService) AddMember(ctx context.Context, requesterUUID, groupUUID, memberUUID uuid.UUID, role string) (*domain.GroupMember, error) {
	if err := validateAssignableRole(role); err != nil {
		return nil, err
	}

	requesterMember, err := s.groupRepo.GetMember(ctx, groupUUID, requesterUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: add member get requester: %w", err)
	}
	if requesterMember == nil || requesterMember.Role != domain.GroupRoleAuthor {
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

func (s *GroupMemberService) UpdateMemberRole(ctx context.Context, requesterUUID, groupUUID, memberUUID uuid.UUID, role string) (*domain.GroupMember, error) {
	if err := validateAssignableRole(role); err != nil {
		return nil, err
	}

	requesterMember, err := s.groupRepo.GetMember(ctx, groupUUID, requesterUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: update member get requester: %w", err)
	}
	if requesterMember == nil || requesterMember.Role != domain.GroupRoleAuthor {
		return nil, domain.ErrForbidden
	}

	member, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: update member get member: %w", err)
	}
	if member == nil {
		return nil, domain.ErrUserNotFound
	}

	if member.Role == domain.GroupRoleAuthor {
		if memberUUID != requesterUUID {
			return nil, domain.ErrForbidden
		}
		if role != domain.GroupRoleAuthor {
			return nil, domain.ErrLastAuthor
		}
		return member, nil
	}

	if err := s.groupRepo.UpdateMemberRole(ctx, groupUUID, memberUUID, role); err != nil {
		return nil, fmt.Errorf("group member service: update member role: %w", err)
	}

	updatedMember, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("group member service: update member fetch updated: %w", err)
	}

	return updatedMember, nil
}

func (s *GroupMemberService) RemoveMember(ctx context.Context, requesterUUID, groupUUID, memberUUID uuid.UUID) error {
	requesterMember, err := s.groupRepo.GetMember(ctx, groupUUID, requesterUUID)
	if err != nil {
		return fmt.Errorf("group member service: remove member get requester: %w", err)
	}
	if requesterMember == nil || requesterMember.Role != domain.GroupRoleAuthor {
		return domain.ErrForbidden
	}

	member, err := s.groupRepo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return fmt.Errorf("group member service: remove member get member: %w", err)
	}
	if member == nil {
		return domain.ErrUserNotFound
	}

	if member.Role == domain.GroupRoleAuthor {
		return domain.ErrLastAuthor
	}

	if err := s.groupRepo.RemoveMember(ctx, groupUUID, memberUUID); err != nil {
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
