package document_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document"
)

// MockGetDocumentService is a mock implementation of the get document service
type MockGetDocumentService struct {
	mock.Mock
}

func (m *MockGetDocumentService) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Document, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

// MockGetAllDocumentsService is a mock implementation of the get all documents service
type MockGetAllDocumentsService struct {
	mock.Mock
}

func (m *MockGetAllDocumentsService) GetAll(ctx context.Context) ([]*domain.Document, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Document), args.Error(1)
}

func TestNewGetDocumentHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SuccessfulGetDocument", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetDocumentService)
		handler := document.NewGetDocumentHandler(mockService)

		documentUUID := uuid.New()
		expectedDocument := &domain.Document{
			UUID:      documentUUID,
			GroupUUID: uuid.New(),
			Name:      "Test Document",
			Content:   "This is test content",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("GetByUUID", mock.Anything, documentUUID).Return(expectedDocument, nil)

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Document", response["name"])
		assert.Equal(t, "This is test content", response["content"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetDocumentService)
		handler := document.NewGetDocumentHandler(mockService)

		req := httptest.NewRequest("GET", "/documents/invalid-uuid", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: "invalid-uuid"}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid uuid format", response["error"])

		mockService.AssertNotCalled(t, "GetByUUID")
	})

	t.Run("DocumentNotFound", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetDocumentService)
		handler := document.NewGetDocumentHandler(mockService)

		documentUUID := uuid.New()
		mockService.On("GetByUUID", mock.Anything, documentUUID).Return(nil, domain.ErrDocumentNotFound)

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "document not found", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetDocumentService)
		handler := document.NewGetDocumentHandler(mockService)

		documentUUID := uuid.New()
		mockService.On("GetByUUID", mock.Anything, documentUUID).Return(nil, domain.ErrInternal)

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get document", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetDocumentService)
		handler := document.NewGetDocumentHandler(mockService)

		documentUUID := uuid.New()
		mockService.On("GetByUUID", mock.Anything, documentUUID).Return(nil, errors.New("database connection failed"))

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get document", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestNewGetAllDocumentsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SuccessfulGetAllDocuments", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetAllDocumentsService)
		handler := document.NewGetAllDocumentsHandler(mockService)

		document1 := &domain.Document{
			UUID:      uuid.New(),
			GroupUUID: uuid.New(),
			Name:      "Document 1",
			Content:   "Content 1",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		document2 := &domain.Document{
			UUID:      uuid.New(),
			GroupUUID: uuid.New(),
			Name:      "Document 2",
			Content:   "Content 2",
			CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		}

		expectedDocuments := []*domain.Document{document1, document2}
		mockService.On("GetAll", mock.Anything).Return(expectedDocuments, nil)

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "documents")

		documents := response["documents"].([]interface{})
		assert.Len(t, documents, 2)
		assert.Equal(t, "Document 1", documents[0].(map[string]interface{})["name"])
		assert.Equal(t, "Document 2", documents[1].(map[string]interface{})["name"])

		mockService.AssertExpectations(t)
	})

	t.Run("EmptyDocumentsList", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetAllDocumentsService)
		handler := document.NewGetAllDocumentsHandler(mockService)

		expectedDocuments := []*domain.Document{}
		mockService.On("GetAll", mock.Anything).Return(expectedDocuments, nil)

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "documents")

		documents := response["documents"].([]interface{})
		assert.Len(t, documents, 0)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetAllDocumentsService)
		handler := document.NewGetAllDocumentsHandler(mockService)

		mockService.On("GetAll", mock.Anything).Return(nil, domain.ErrInternal)

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get documents", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService := new(MockGetAllDocumentsService)
		handler := document.NewGetAllDocumentsHandler(mockService)

		mockService.On("GetAll", mock.Anything).Return(nil, errors.New("database connection failed"))

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to get documents", response["error"])

		mockService.AssertExpectations(t)
	})
}
