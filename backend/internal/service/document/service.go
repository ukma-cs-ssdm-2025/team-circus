package document

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/document"
)

type DocumentService struct {
	repo *document.DocumentRepository
}

func NewDocumentService(repo *document.DocumentRepository) *DocumentService {
	return &DocumentService{
		repo: repo,
	}
}
