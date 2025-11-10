package document_test

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/testutil"
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

const documentsEndpoint = "/documents"

type documentContextBuilder func(t *testing.T, documentUUID, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder)

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

		c, w := authedDocumentContext(t, documentUUID, userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "Test Document", response["name"])
		assert.Equal(t, "This is test content", response["content"])

		mockService.AssertExpectations(t)
	})

	main.Run("FailureCases", func(t *testing.T) {
		cases := []struct {
			name          string
			setupMock     func(*mockGetDocumentService, uuid.UUID, uuid.UUID)
			buildContext  documentContextBuilder
			expectedCode  int
			expectedError string
			expectCall    bool
		}{
			{
				name: "InvalidUUID",
				buildContext: func(t *testing.T, _ uuid.UUID, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
					return documentContextWithParams(t, "invalid-uuid", &userUUID)
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "invalid uuid format",
			},
			{
				name: "MissingUserContext",
				buildContext: func(t *testing.T, documentUUID, _ uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
					return documentContextWithParams(t, documentUUID.String(), nil)
				},
				expectedCode:  http.StatusUnauthorized,
				expectedError: "user context missing",
			},
			{
				name: "DocumentNotFound",
				setupMock: func(m *mockGetDocumentService, documentUUID, userUUID uuid.UUID) {
					m.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, domain.ErrDocumentNotFound)
				},
				buildContext: func(t *testing.T, documentUUID, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
					return authedDocumentContext(t, documentUUID, userUUID)
				},
				expectedCode:  http.StatusNotFound,
				expectedError: "document not found",
				expectCall:    true,
			},
			{
				name: "ServiceInternalError",
				setupMock: func(m *mockGetDocumentService, documentUUID, userUUID uuid.UUID) {
					m.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, domain.ErrInternal)
				},
				buildContext: func(t *testing.T, documentUUID, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
					return authedDocumentContext(t, documentUUID, userUUID)
				},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to get document",
				expectCall:    true,
			},
			{
				name: "ServiceGenericError",
				setupMock: func(m *mockGetDocumentService, documentUUID, userUUID uuid.UUID) {
					m.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, errors.New("database connection failed"))
				},
				buildContext: func(t *testing.T, documentUUID, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
					return authedDocumentContext(t, documentUUID, userUUID)
				},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to get document",
				expectCall:    true,
			},
			{
				name: "Forbidden",
				setupMock: func(m *mockGetDocumentService, documentUUID, userUUID uuid.UUID) {
					m.On("GetByUUIDForUser", mock.Anything, documentUUID, userUUID).Return(nil, domain.ErrForbidden)
				},
				buildContext: func(t *testing.T, documentUUID, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
					return authedDocumentContext(t, documentUUID, userUUID)
				},
				expectedCode:  http.StatusForbidden,
				expectedError: "access forbidden",
				expectCall:    true,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				mockService, handler := setup(t)
				documentUUID := uuid.New()
				userUUID := uuid.New()

				if tc.setupMock != nil {
					tc.setupMock(mockService, documentUUID, userUUID)
				}

				c, w := tc.buildContext(t, documentUUID, userUUID)

				handler(c)

				assert.Equal(t, tc.expectedCode, w.Code)

				response := testutil.DecodeResponse(t, w)
				assert.Equal(t, tc.expectedError, response["error"])

				if tc.expectCall {
					mockService.AssertExpectations(t)
				} else {
					mockService.AssertNotCalled(t, "GetByUUIDForUser")
				}
			})
		}
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

		c, w := authedDocumentsListContext(t, userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		response := testutil.DecodeResponse(t, w)
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

		c, w := authedDocumentsListContext(t, userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		response := testutil.DecodeResponse(t, w)
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

		c, w := authedDocumentsListContext(t, userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "failed to get documents", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceGenericError", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		userUUID := uuid.New()
		mockService.On("GetAllForUser", mock.Anything, userUUID).Return(nil, errors.New("database connection failed"))

		c, w := authedDocumentsListContext(t, userUUID)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "failed to get documents", response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("MissingUserContext", func(t *testing.T) {
		mockService, handler := setup(t)

		c, w := testutil.NewRequestContext(t, http.MethodGet, documentsEndpoint)

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertNotCalled(t, "GetAllForUser")
	})
}

func documentContextWithParams(t *testing.T, uuidValue string, userUUID *uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
	path := fmt.Sprintf("%s/%s", documentsEndpoint, uuidValue)
	c, w := testutil.NewRequestContext(t, http.MethodGet, path)
	c.Params = gin.Params{{Key: "uuid", Value: uuidValue}}
	if userUUID != nil {
		c.Set("user_uid", *userUUID)
	}
	return c, w
}

func authedDocumentContext(t *testing.T, documentUUID, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
	return documentContextWithParams(t, documentUUID.String(), &userUUID)
}

func authedDocumentsListContext(t *testing.T, userUUID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
	c, w := testutil.NewRequestContext(t, http.MethodGet, documentsEndpoint)
	c.Set("user_uid", userUUID)
	return c, w
}
