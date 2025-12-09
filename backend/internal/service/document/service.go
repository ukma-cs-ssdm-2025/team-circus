package document

import (
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
