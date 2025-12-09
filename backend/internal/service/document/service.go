package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/document"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/member"
)

type DocumentService struct {
	repo       *document.DocumentRepository
	memberRepo *member.MemberRepository
	shareCfg   ShareConfig
}

type ShareConfig struct {
	Secret                string
	BaseURL               string
	DefaultExpirationDays int
	MaxExpirationDays     int
}

func NewDocumentService(
	repo *document.DocumentRepository,
	memberRepo *member.MemberRepository,
	shareCfg ShareConfig,
) *DocumentService {
	return &DocumentService{
		repo:       repo,
		memberRepo: memberRepo,
		shareCfg:   shareCfg,
	}
}

// GetMemberRole returns the role of a user within the document's group context.
func (s *DocumentService) GetMemberRole(ctx context.Context, documentUUID, userUUID uuid.UUID) (string, error) {
	doc, err := s.GetByUUID(ctx, documentUUID)
	if err != nil {
		return "", err
	}

	member, err := s.memberRepo.GetMember(ctx, doc.GroupUUID, userUUID)
	if err != nil {
		return "", fmt.Errorf("document service: getMemberRole: %w", err)
	}
	if member == nil {
		return "", domain.ErrForbidden
	}

	return member.Role, nil
}
