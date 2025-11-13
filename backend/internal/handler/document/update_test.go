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

type mockUpdateDocumentService struct {
	mock.Mock
}

func (m *mockUpdateDocumentService) Update(ctx context.Context, uuid uuid.UUID, name, content string) (*domain.Document, error) {
	args := m.Called(ctx, uuid, name, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1) //nolint:errcheck
}

func TestNewUpdateDocumentHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockUpdateDocumentService, gin.HandlerFunc) {
		mockService := &mockUpdateDocumentService{}
		handler := document.NewUpdateDocumentHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulUpdate", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		expectedDocument := &domain.Document{
			UUID:      documentUUID,
			GroupUUID: uuid.New(),
			Name:      "Updated Document",
			Content:   "Updated content",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("Update", mock.Anything, documentUUID, "Updated Document", "Updated content").Return(expectedDocument, nil)

		requestBody := requests.UpdateDocumentRequest{
			Name:    "Updated Document",
			Content: "Updated content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/documents/"+documentUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Document", response["name"])
		assert.Equal(t, "Updated content", response["content"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		requestBody := requests.UpdateDocumentRequest{
			Name:    "Updated Document",
			Content: "Updated content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/documents/invalid-uuid", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: "invalid-uuid"}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid uuid format", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		invalidJSON := `{"name": "Updated Document", "content": "Updated content"` // Missing closing brace

		req := httptest.NewRequest("PUT", "/documents/"+documentUUID.String(), bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request format", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("ValidationFailed_EmptyName", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		requestBody := requests.UpdateDocumentRequest{
			Name:    "",
			Content: "Updated content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/documents/"+documentUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "validation failed", response["error"])

		mockService.AssertNotCalled(t, "Update")
	})

	t.Run("DocumentNotFound", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		mockService.On("Update", mock.Anything, documentUUID, "Updated Document", "Updated content").Return(nil, domain.ErrDocumentNotFound)

		requestBody := requests.UpdateDocumentRequest{
			Name:    "Updated Document",
			Content: "Updated content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/documents/"+documentUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "document not found", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		mockService.On("Update", mock.Anything, documentUUID, "Updated Document", "Updated content").Return(nil, domain.ErrInternal)

		requestBody := requests.UpdateDocumentRequest{
			Name:    "Updated Document",
			Content: "Updated content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/documents/"+documentUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to update document", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		mockService.On(
			"Update",
			mock.Anything,
			documentUUID,
			"Updated Document",
			"Updated content",
		).Return(nil, errors.New("database connection failed"))

		requestBody := requests.UpdateDocumentRequest{
			Name:    "Updated Document",
			Content: "Updated content",
		}

		jsonBody, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/documents/"+documentUUID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "uuid", Value: documentUUID.String()}}

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "failed to update document", response["error"])

		mockService.AssertExpectations(t)
	})
}
