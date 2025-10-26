package document

import (
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
