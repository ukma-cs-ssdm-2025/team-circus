package reg_test

import (
	"context"
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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/testutil"
	"go.uber.org/zap"
)

// MockRegService is a mock implementation of the registration service
type MockRegService struct {
	mock.Mock
}

const registerEndpoint = "/signup"

func (m *MockRegService) Register(ctx context.Context, login string, email string, password string) (*domain.User, error) {
	args := m.Called(ctx, login, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck
}

func TestNewRegHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*MockRegService, gin.HandlerFunc) {
		mockService := &MockRegService{}
		handler := reg.NewRegHandler(mockService, zap.NewNop())
		t.Cleanup(func() {
			mockService.AssertExpectations(t)
		})
		return mockService, handler
	}
	t.Run("SuccessfulRegistration", func(t *testing.T) {
		// Arrange
		mockService, handler := setup(t)

		expectedUser := &domain.User{
			UUID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:     "testuser",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockService.On("Register", mock.Anything, "testuser", "test@example.com", "testpassword123").
			Return(expectedUser, nil)

		requestBody := requests.RegRequest{
			Login:    "testuser",
			Email:    "test@example.com",
			Password: "testpassword123",
		}

		c, w := testutil.NewJSONContext(t, http.MethodPost, registerEndpoint, requestBody)

		// Act
		handler(c)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		response := testutil.DecodeResponse(t, w)
		assert.Equal(t, "testuser", response["login"])
		assert.Equal(t, "test@example.com", response["email"])

		mockService.AssertExpectations(t)
	})

	t.Run("FailureCases", func(t *testing.T) {
		type (
			registerSetup      func(*MockRegService)
			registerContextGen func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder)
		)

		cases := []struct {
			name          string
			setup         registerSetup
			buildContext  registerContextGen
			expectedCode  int
			expectedError string
			expectCall    bool
		}{
			{
				name: "InvalidJSON",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					payload := `{"login": "testuser", "email": "test@example.com", "password": "testpassword123"`
					return testutil.NewRawContext(t, http.MethodPost, registerEndpoint, []byte(payload), "application/json")
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "invalid request format",
			},
			{
				name: "ValidationFailure",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					body := requests.RegRequest{Login: "", Email: "test@example.com", Password: "testpassword123"}
					return testutil.NewJSONContext(t, http.MethodPost, registerEndpoint, body)
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "validation failed",
			},
			{
				name: "ServiceInternalError",
				setup: func(service *MockRegService) {
					service.On("Register", mock.Anything, "testuser", "test@example.com", "testpassword123").
						Return(nil, domain.ErrInternal)
				},
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					body := requests.RegRequest{Login: "testuser", Email: "test@example.com", Password: "testpassword123"}
					return testutil.NewJSONContext(t, http.MethodPost, registerEndpoint, body)
				},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to register",
				expectCall:    true,
			},
			{
				name: "ServiceGenericError",
				setup: func(service *MockRegService) {
					service.On("Register", mock.Anything, "testuser", "test@example.com", "testpassword123").
						Return(nil, errors.New("database connection failed"))
				},
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					body := requests.RegRequest{Login: "testuser", Email: "test@example.com", Password: "testpassword123"}
					return testutil.NewJSONContext(t, http.MethodPost, registerEndpoint, body)
				},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to register",
				expectCall:    true,
			},
			{
				name: "MissingContentType",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					return testutil.NewRawContext(t, http.MethodPost, registerEndpoint, []byte("invalid json"), "")
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "invalid request format",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				mockService, handler := setup(t)
				if tc.setup != nil {
					tc.setup(mockService)
				}

				c, w := tc.buildContext(t)

				handler(c)

				assert.Equal(t, tc.expectedCode, w.Code)

				response := testutil.DecodeResponse(t, w)
				assert.Equal(t, tc.expectedError, response["error"])

				if tc.expectCall {
					mockService.AssertExpectations(t)
				} else {
					mockService.AssertNotCalled(t, "Register")
				}
			})
		}
	})
}
