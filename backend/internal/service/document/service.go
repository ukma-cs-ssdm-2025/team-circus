package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/document"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
)

type DocumentService struct {
	repo      *document.DocumentRepository
	groupRepo *grouprepo.GroupRepository
}

func NewDocumentService(repo *document.DocumentRepository, groupRepo *grouprepo.GroupRepository) *DocumentService {
	return &DocumentService{
		repo:      repo,
		groupRepo: groupRepo,
	}
}

func (s *DocumentService) ensureCanEditDocuments(ctx context.Context, groupUUID, userUUID uuid.UUID) error {
	member, err := s.groupRepo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return fmt.Errorf("document service: ensureCanEditDocuments get member: %w", err)
	}

	if member == nil {
		return domain.ErrForbidden
	}

	if member.Role == domain.GroupRoleReviewer {
		return domain.ErrForbidden
	}

	return nil
}
