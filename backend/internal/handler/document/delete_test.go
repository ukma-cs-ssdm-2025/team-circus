package document_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document"
	"go.uber.org/zap"
)

type mockDeleteDocumentService struct {
	mock.Mock
}

func (m *mockDeleteDocumentService) Delete(ctx context.Context, uuid uuid.UUID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func TestNewDeleteDocumentHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*mockDeleteDocumentService, gin.HandlerFunc) {
		mockService := &mockDeleteDocumentService{}
		handler := document.NewDeleteDocumentHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}

	t.Run("SuccessfulDelete", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		mockService.On("Delete", mock.Anything, documentUUID).Return(nil)

		req := httptest.NewRequest("DELETE", "/documents/"+documentUUID.String(), nil)
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
		assert.Equal(t, "Document deleted successfully", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		req := httptest.NewRequest("DELETE", "/documents/invalid-uuid", nil)
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

		mockService.AssertNotCalled(t, "Delete")
	})

	t.Run("DocumentNotFound", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		mockService.On("Delete", mock.Anything, documentUUID).Return(domain.ErrDocumentNotFound)

		req := httptest.NewRequest("DELETE", "/documents/"+documentUUID.String(), nil)
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
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		mockService.On("Delete", mock.Anything, documentUUID).Return(domain.ErrInternal)

		req := httptest.NewRequest("DELETE", "/documents/"+documentUUID.String(), nil)
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
		assert.Equal(t, "failed to delete document", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		documentUUID := uuid.New()
		mockService.On("Delete", mock.Anything, documentUUID).Return(errors.New("database connection failed"))

		req := httptest.NewRequest("DELETE", "/documents/"+documentUUID.String(), nil)
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
		assert.Equal(t, "failed to delete document", response["error"])

		mockService.AssertExpectations(t)
	})
}

