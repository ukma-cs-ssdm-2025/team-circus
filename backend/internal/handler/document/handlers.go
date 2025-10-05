package document

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/responses"
)

func mapDocumentToCreateResponse(document *domain.Document) responses.CreateDocumentResponse {
	return responses.CreateDocumentResponse{
		UUID:      document.UUID,
		GroupUUID: document.GroupUUID,
		Name:      document.Name,
		Content:   document.Content,
		CreatedAt: document.CreatedAt,
	}
}

func mapDocumentToGetResponse(document *domain.Document) responses.GetDocumentResponse {
	return responses.GetDocumentResponse{
		UUID:      document.UUID,
		GroupUUID: document.GroupUUID,
		Name:      document.Name,
		Content:   document.Content,
		CreatedAt: document.CreatedAt,
	}
}

func mapDocumentToUpdateResponse(document *domain.Document) responses.UpdateDocumentResponse {
	return responses.UpdateDocumentResponse{
		UUID:      document.UUID,
		GroupUUID: document.GroupUUID,
		Name:      document.Name,
		Content:   document.Content,
		CreatedAt: document.CreatedAt,
	}
}

func mapDocumentsToGetAllResponse(documents []*domain.Document) []responses.GetDocumentResponse {
	result := make([]responses.GetDocumentResponse, len(documents))
	for i, document := range documents {
		result[i] = mapDocumentToGetResponse(document)
	}
	return result
}
