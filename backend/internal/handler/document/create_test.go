package document_test

import (
	"bytes"
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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document/requests"
	"go.uber.org/zap"
)

type mockCreateDocumentService struct {
	mock.Mock
}

func (m *mockCreateDocumentService) Create(
	ctx context.Context,
	userUUID, groupUUID uuid.UUID,
	name, content string,
) (*domain.Document, error) {
	args := m.Called(ctx, userUUID, groupUUID, name, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1) //nolint:errcheck
}

func TestNewCreateDocumentHandler(main *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockCreateDocumentService, gin.HandlerFunc) {
		mockService := &mockCreateDocumentService{}
		handler := document.NewCreateDocumentHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	main.Run("SuccessfulCreateDocument", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		expectedDocument := &domain.Document{
			UUID:      uuid.New(),
			GroupUUID: groupUUID,
			Name:      "Test Document",
			Content:   "This is test content",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("Create", mock.Anything, userUUID, groupUUID, "Test Document", "This is test content").
			Return(expectedDocument, nil)

		requestBody := requests.CreateDocumentRequest{
			GroupUUID: groupUUID,
			Name:      "Test Document",
			Content:   "This is test content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/documents", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Document", response["name"])
		assert.Equal(t, "This is test content", response["content"])

		mockService.AssertExpectations(t)
	})

	main.Run("MissingUserContext", func(t *testing.T) {
		mockService, handler := setup(t)

		requestBody := requests.CreateDocumentRequest{
			GroupUUID: uuid.New(),
			Name:      "Test Document",
			Content:   "This is test content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/documents", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	main.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		invalidJSON := `{"group_uuid": "123e4567-e89b-12d3-a456-426614174000", "name": "Test Document", "content": "This is test content"` //nolint:revive
		userUUID := uuid.New()

		req := httptest.NewRequest("POST", "/documents", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request format", response["error"])

		mockService.AssertNotCalled(t, "Create")
	})

	main.Run("ValidationFailure", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		requestBody := requests.CreateDocumentRequest{
			GroupUUID: uuid.New(),
			Name:      "", // Empty name should fail validation
			Content:   "This is test content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/documents", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation failed", response["error"])
		assert.Contains(t, response, "details")

		mockService.AssertNotCalled(t, "Create")
	})

	main.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("Create", mock.Anything, userUUID, groupUUID, "Test Document", "This is test content").
			Return(nil, domain.ErrInternal)

		requestBody := requests.CreateDocumentRequest{
			GroupUUID: groupUUID,
			Name:      "Test Document",
			Content:   "This is test content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/documents", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to create document", response["error"])

		mockService.AssertExpectations(t)
	})

	main.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		groupUUID := uuid.New()
		userUUID := uuid.New()
		mockService.On("Create", mock.Anything, userUUID, groupUUID, "Test Document", "This is test content").
			Return(nil, errors.New("database connection failed"))

		requestBody := requests.CreateDocumentRequest{
			GroupUUID: groupUUID,
			Name:      "Test Document",
			Content:   "This is test content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/documents", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to create document", response["error"])

		mockService.AssertExpectations(t)
	})

	main.Run("MissingContentType", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()

		// Send invalid JSON without Content-Type
		req := httptest.NewRequest("POST", "/documents", bytes.NewBufferString("invalid json"))
		// Don't set Content-Type header
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_uid", userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request format", response["error"])

		mockService.AssertNotCalled(t, "Create")
	})
}
