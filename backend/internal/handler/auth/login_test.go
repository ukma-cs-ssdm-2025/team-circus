package auth_test

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
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/testutil"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	args := m.Called(ctx, login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck
}

func (m *MockUserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck
}

const loginEndpoint = "/auth/login"

func TestNewLogInHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	setup := func(t *testing.T) (*MockUserRepository, gin.HandlerFunc) {
		mockRepo := new(MockUserRepository)
		handler := auth.NewLogInHandler(mockRepo, zap.NewNop(), "test-secret-key", 10, 4320)
		t.Cleanup(func() {
			mockRepo.AssertExpectations(t)
		})
		return mockRepo, handler
	}

	t.Run("SuccessfulLogin", func(t *testing.T) {
		mockRepo, handler := setup(t)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword123"), bcrypt.DefaultCost)
		assert.NoError(t, err)
		expectedUser := &domain.User{
			UUID:      uuid.New(),
			Login:     "testuser",
			Email:     "test@example.com",
			Password:  string(hashedPassword),
			CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockRepo.On("GetByLogin", mock.Anything, "testuser").Return(expectedUser, nil)

		body := requests.LogInRequest{Login: "testuser", Password: "testpassword123"}
		c, w := testutil.NewJSONContext(t, http.MethodPost, loginEndpoint, body)

		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 2)
		assert.NotNil(t, testutil.CookieByName(cookies, "accessToken"))
		assert.NotNil(t, testutil.CookieByName(cookies, "refreshToken"))
	})

	t.Run("FailureCases", func(t *testing.T) {
		type (
			loginSetup      func(*MockUserRepository)
			loginContextGen func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder)
		)

		cases := []struct {
			name          string
			setup         loginSetup
			buildContext  loginContextGen
			expectedCode  int
			expectedError string
			expectCall    bool
		}{
			{
				name: "InvalidJSON",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					payload := `{"login": "testuser", "password": "testpassword123"`
					return testutil.NewRawContext(t, http.MethodPost, loginEndpoint, []byte(payload), "application/json")
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "invalid request format",
			},
			{
				name: "UserNotFound",
				setup: func(repo *MockUserRepository) {
					repo.On("GetByLogin", mock.Anything, "nonexistent").Return(nil, nil)
				},
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					body := requests.LogInRequest{Login: "nonexistent", Password: "testpassword123"}
					return testutil.NewJSONContext(t, http.MethodPost, loginEndpoint, body)
				},
				expectedCode:  http.StatusUnauthorized,
				expectedError: "invalid credentials",
				expectCall:    true,
			},
			{
				name: "WrongPassword",
				setup: func(repo *MockUserRepository) {
					hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
					repo.On("GetByLogin", mock.Anything, "testuser").Return(&domain.User{
						UUID:      uuid.New(),
						Login:     "testuser",
						Email:     "test@example.com",
						Password:  string(hashed),
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					}, nil)
				},
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					body := requests.LogInRequest{Login: "testuser", Password: "wrongpassword"}
					return testutil.NewJSONContext(t, http.MethodPost, loginEndpoint, body)
				},
				expectedCode:  http.StatusUnauthorized,
				expectedError: "invalid credentials",
				expectCall:    true,
			},
			{
				name: "ServiceInternalError",
				setup: func(repo *MockUserRepository) {
					repo.On("GetByLogin", mock.Anything, "testuser").Return(nil, domain.ErrInternal)
				},
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					body := requests.LogInRequest{Login: "testuser", Password: "testpassword123"}
					return testutil.NewJSONContext(t, http.MethodPost, loginEndpoint, body)
				},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to log in",
				expectCall:    true,
			},
			{
				name: "ServiceGenericError",
				setup: func(repo *MockUserRepository) {
					repo.On("GetByLogin", mock.Anything, "testuser").Return(nil, errors.New("database connection failed"))
				},
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					body := requests.LogInRequest{Login: "testuser", Password: "testpassword123"}
					return testutil.NewJSONContext(t, http.MethodPost, loginEndpoint, body)
				},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to log in",
				expectCall:    true,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				mockRepo, handler := setup(t)
				if tc.setup != nil {
					tc.setup(mockRepo)
				}

				c, w := tc.buildContext(t)

				handler(c)

				assert.Equal(t, tc.expectedCode, w.Code)

				response := testutil.DecodeResponse(t, w)
				assert.Equal(t, tc.expectedError, response["error"])

				if tc.expectCall {
					mockRepo.AssertExpectations(t)
				} else {
					mockRepo.AssertNotCalled(t, "GetByLogin")
				}
			})
		}
	})
}
