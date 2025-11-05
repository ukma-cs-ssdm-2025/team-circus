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
	userUUID,
	groupUUID uuid.UUID,
	name,
	content string,
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

	successfulUserUUID := uuid.New()
	successfulGroupUUID := uuid.New()
	expectedDocument := &domain.Document{
		UUID:      uuid.New(),
		GroupUUID: successfulGroupUUID,
		Name:      "Test Document",
		Content:   "This is test content",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	forbiddenUserUUID := uuid.New()
	forbiddenGroupUUID := uuid.New()

	internalErrorUserUUID := uuid.New()
	internalErrorGroupUUID := uuid.New()

	genericErrorUserUUID := uuid.New()
	genericErrorGroupUUID := uuid.New()

	validationUserUUID := uuid.New()
	invalidJSONUserUUID := uuid.New()
	missingContentTypeUserUUID := uuid.New()

	testCases := []struct {
		name           string
		body           interface{}
		rawBody        string
		contentType    string
		userUUID       uuid.UUID
		expectedStatus int
		setupMock      func(*mockCreateDocumentService)
		assertResponse func(*testing.T, map[string]interface{})
		assertService  func(*testing.T, *mockCreateDocumentService)
	}{
		{
			name:           "SuccessfulCreateDocument",
			body:           requests.CreateDocumentRequest{GroupUUID: successfulGroupUUID, Name: "Test Document", Content: "This is test content"},
			contentType:    "application/json",
			userUUID:       successfulUserUUID,
			expectedStatus: http.StatusCreated,
			setupMock: func(mockService *mockCreateDocumentService) {
				mockService.On("Create", mock.Anything, successfulUserUUID, successfulGroupUUID, "Test Document", "This is test content").
					Return(expectedDocument, nil)
			},
			assertResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "Test Document", response["name"])
				assert.Equal(t, "This is test content", response["content"])
			},
		},
		{
			name:           "InvalidJSON",
			rawBody:        `{"group_uuid": "123e4567-e89b-12d3-a456-426614174000", "name": "Test Document", "content": "This is test content"`,
			contentType:    "application/json",
			userUUID:       invalidJSONUserUUID,
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "invalid request format", response["error"])
			},
			assertService: func(t *testing.T, mockService *mockCreateDocumentService) {
				mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:           "ValidationFailure",
			body:           requests.CreateDocumentRequest{GroupUUID: uuid.New(), Name: "", Content: "This is test content"},
			contentType:    "application/json",
			userUUID:       validationUserUUID,
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "validation failed", response["error"])
				assert.Contains(t, response, "details")
			},
			assertService: func(t *testing.T, mockService *mockCreateDocumentService) {
				mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
		},
		{
			name:           "ForbiddenResponse",
			body:           requests.CreateDocumentRequest{GroupUUID: forbiddenGroupUUID, Name: "Test Document", Content: "content"},
			contentType:    "application/json",
			userUUID:       forbiddenUserUUID,
			expectedStatus: http.StatusForbidden,
			setupMock: func(mockService *mockCreateDocumentService) {
				mockService.On("Create", mock.Anything, forbiddenUserUUID, forbiddenGroupUUID, "Test Document", "content").
					Return(nil, domain.ErrForbidden)
			},
			assertResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "access forbidden", response["error"])
			},
		},
		{
			name:           "ServiceInternalError",
			body:           requests.CreateDocumentRequest{GroupUUID: internalErrorGroupUUID, Name: "Test Document", Content: "This is test content"},
			contentType:    "application/json",
			userUUID:       internalErrorUserUUID,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mockService *mockCreateDocumentService) {
				mockService.On("Create", mock.Anything, internalErrorUserUUID, internalErrorGroupUUID, "Test Document", "This is test content").
					Return(nil, domain.ErrInternal)
			},
			assertResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "failed to create document", response["error"])
			},
		},
		{
			name:           "ServiceGenericError",
			body:           requests.CreateDocumentRequest{GroupUUID: genericErrorGroupUUID, Name: "Test Document", Content: "This is test content"},
			contentType:    "application/json",
			userUUID:       genericErrorUserUUID,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mockService *mockCreateDocumentService) {
				mockService.On("Create", mock.Anything, genericErrorUserUUID, genericErrorGroupUUID, "Test Document", "This is test content").
					Return(nil, errors.New("database connection failed"))
			},
			assertResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "failed to create document", response["error"])
			},
		},
		{
			name:           "MissingContentType",
			rawBody:        "invalid json",
			userUUID:       missingContentTypeUserUUID,
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "invalid request format", response["error"])
			},
			assertService: func(t *testing.T, mockService *mockCreateDocumentService) {
				mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
		},
	}

	for _, tc := range testCases {
		main.Run(tc.name, func(t *testing.T) {
			mockService, handler := setup(t)
			if tc.setupMock != nil {
				tc.setupMock(mockService)
			}

			bodyReader := prepareRequestBody(t, tc.body, tc.rawBody)
			recorder := performCreateDocumentRequest(t, handler, tc.userUUID, tc.contentType, bodyReader)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			if tc.assertResponse != nil {
				response := unmarshalResponse(t, recorder)
				tc.assertResponse(t, response)
			}

			if tc.assertService != nil {
				tc.assertService(t, mockService)
			}
		})
	}
}

func prepareRequestBody(t *testing.T, body interface{}, rawBody string) *bytes.Buffer {
	t.Helper()

	switch {
	case rawBody != "":
		return bytes.NewBufferString(rawBody)
	case body != nil:
		jsonBody, err := json.Marshal(body)
		assert.NoError(t, err)
		return bytes.NewBuffer(jsonBody)
	default:
		return bytes.NewBuffer(nil)
	}
}

func performCreateDocumentRequest(
	t *testing.T,
	handler gin.HandlerFunc,
	userUUID uuid.UUID,
	contentType string,
	body *bytes.Buffer,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/documents", body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_uid", userUUID)

	handler(c)

	return w
}

func unmarshalResponse(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	return response
}
