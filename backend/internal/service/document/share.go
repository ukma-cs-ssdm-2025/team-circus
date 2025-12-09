package document

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

const (
	minShareExpirationDays = 1
	defaultMaxExpiration   = 90
)

var ErrInvalidExpiration = errors.New("invalid expiration days")

func (s *DocumentService) GenerateShareLink(
	ctx context.Context,
	docUUID, userUUID uuid.UUID,
	expirationDays int,
) (string, time.Time, error) {
	if s.shareCfg.Secret == "" || s.shareCfg.BaseURL == "" {
		return "", time.Time{}, domain.ErrInternal
	}

	doc, err := s.GetByUUID(ctx, docUUID)
	if err != nil {
		return "", time.Time{}, err
	}

	member, err := s.memberRepo.GetMember(ctx, doc.GroupUUID, userUUID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("document service: share: %w", err)
	}
	if member == nil || member.Role == domain.RoleViewer {
		return "", time.Time{}, domain.ErrForbidden
	}

	effectiveDays := expirationDays
	if effectiveDays == 0 {
		effectiveDays = s.shareCfg.DefaultExpirationDays
	}
	if effectiveDays == 0 {
		effectiveDays = minShareExpirationDays
	}

	maxDays := s.shareCfg.MaxExpirationDays
	if maxDays <= 0 {
		maxDays = defaultMaxExpiration
	}

	if effectiveDays < minShareExpirationDays || effectiveDays > maxDays {
		return "", time.Time{}, ErrInvalidExpiration
	}

	expiresAt := time.Now().UTC().Add(time.Duration(effectiveDays) * 24 * time.Hour).Truncate(time.Second)

	signature := s.computeSignature(docUUID, expiresAt)
	shareURL := fmt.Sprintf(
		"%s/documents/public?doc=%s&sig=%s&exp=%d",
		strings.TrimRight(s.shareCfg.BaseURL, "/"),
		docUUID.String(),
		signature,
		expiresAt.Unix(),
	)

	return shareURL, expiresAt, nil
}

func (s *DocumentService) ValidateShareLink(
	ctx context.Context,
	docParam, sigParam, expParam string,
) (*domain.Document, error) {
	if s.shareCfg.Secret == "" {
		return nil, domain.ErrInternal
	}

	if docParam == "" || sigParam == "" || expParam == "" {
		return nil, domain.ErrShareLinkInvalid
	}

	docUUID, err := uuid.Parse(docParam)
	if err != nil {
		return nil, domain.ErrShareLinkInvalid
	}

	expTS, err := strconv.ParseInt(expParam, 10, 64)
	if err != nil {
		return nil, domain.ErrShareLinkInvalid
	}

	expiresAt := time.Unix(expTS, 0).UTC()
	if time.Now().UTC().After(expiresAt) {
		return nil, domain.ErrShareLinkExpired
	}

	expectedSig := s.computeSignature(docUUID, expiresAt)
	if !hmac.Equal([]byte(expectedSig), []byte(sigParam)) {
		return nil, domain.ErrShareLinkInvalid
	}

	doc, err := s.GetByUUID(ctx, docUUID)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) computeSignature(docUUID uuid.UUID, expiresAt time.Time) string {
	mac := hmac.New(sha256.New, []byte(s.shareCfg.Secret))
	payload := fmt.Sprintf("%s:%d", docUUID.String(), expiresAt.Unix())
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
