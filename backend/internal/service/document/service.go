package document

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/document"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/member"
)

type DocumentService struct {
	repo       *document.DocumentRepository
	memberRepo *member.MemberRepository
}

func NewDocumentService(repo *document.DocumentRepository, memberRepo *member.MemberRepository) *DocumentService {
	return &DocumentService{
		repo:       repo,
		memberRepo: memberRepo,
	}
}
