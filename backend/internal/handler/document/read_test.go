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
	"go.uber.org/zap"
)

type mockGetDocumentService struct {
	mock.Mock
}

func (m *mockGetDocumentService) GetByUUIDForUser(ctx context.Context, documentUUID, userUUID uuid.UUID) (*domain.Document, error) {
	args := m.Called(ctx, documentUUID, userUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1) //nolint:errcheck
}

type mockGetAllDocumentsService struct {
	mock.Mock
}

func (m *mockGetAllDocumentsService) GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Document, error) {
	args := m.Called(ctx, userUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Document), args.Error(1) //nolint:errcheck
}

func TestNewGetDocumentHandler(main *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockGetDocumentService, gin.HandlerFunc) {
		mockService := &mockGetDocumentService{}
		handler := document.NewGetDocumentHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	main.Run("SuccessfulGetDocument", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		userUUID := uuid.New()
		expectedDocument := &domain.Document{
			UUID:      documentUUID,
			GroupUUID: uuid.New(),
			Name:      "Test Document",
			Content:   "This is test content",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(expectedDocument, nil)

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}
		c.Set("user_uid", userUUID)

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

	main.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()

		req := httptest.NewRequest("GET", "/documents/invalid-uuid", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: "invalid-uuid"}}
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid uuid format", response["error"])

		mockService.AssertNotCalled(t, "GetByUUIDForUser")
	})

	main.Run("MissingUserContext", func(t *testing.T) {
		mockService, handler := setup(t)

		documentUUID := uuid.New()

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertNotCalled(t, "GetByUUIDForUser")
	})

	main.Run("DocumentNotFound", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, domain.ErrDocumentNotFound)

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}
		c.Set("user_uid", userUUID)

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

	main.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, domain.ErrInternal)

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}
		c.Set("user_uid", userUUID)

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

	main.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, errors.New("database connection failed"))

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}
		c.Set("user_uid", userUUID)

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

	main.Run("Forbidden", func(t *testing.T) {
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, domain.ErrForbidden)

		req := httptest.NewRequest("GET", "/documents/"+documentUUID.String(), nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}
		c.Set("user_uid", userUUID)

		handler(c)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access forbidden", response["error"])

		mockService.AssertExpectations(t)
	})
}

func TestNewGetAllDocumentsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockGetAllDocumentsService, gin.HandlerFunc) {
		mockService := &mockGetAllDocumentsService{}
		handler := document.NewGetAllDocumentsHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulGetAllDocuments", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()

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
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(expectedDocuments, nil)

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "documents")

		documents := response["documents"].([]interface{}) //nolint:errcheck
		assert.Len(t, documents, 2)
		assert.Equal(t, "Document 1", documents[0].(map[string]interface{})["name"]) //nolint:errcheck
		assert.Equal(t, "Document 2", documents[1].(map[string]interface{})["name"]) //nolint:errcheck

		mockService.AssertExpectations(t)
	})

	t.Run("EmptyDocumentsList", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		expectedDocuments := []*domain.Document{}
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(expectedDocuments, nil)

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "documents")

		documents := response["documents"].([]interface{}) //nolint:errcheck
		assert.Len(t, documents, 0)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(nil, domain.ErrInternal)

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

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
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(nil, errors.New("database connection failed"))

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

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

	t.Run("MissingUserContext", func(t *testing.T) {
		mockService, handler := setup(t)

		req := httptest.NewRequest("GET", "/documents", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertNotCalled(t, "GetAllForUser")
	})
}
