package auth_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth/requests"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/mocks"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/testutil"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	loginEndpoint   = "/auth/login"
	testSecretToken = "test-secret-key"
	testAccessDur   = 10
	testRefreshDur  = 4320
)

var (
	testPassword       = "testpassword123"
	hashedTestPassword = mustHash(testPassword)
)

func TestNewLogInHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	buildRouter := func(repo *mocks.MockRepository) *gin.Engine {
		router := gin.New()
		router.POST(loginEndpoint, auth.NewLogInHandler(repo, zap.NewNop(), testSecretToken, testAccessDur, testRefreshDur))
		return router
	}

	t.Run("SuccessfulLogin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mocks.NewMockRepository(ctrl)
		repo.EXPECT().
			GetByLogin(gomock.Any(), "testuser").
			Return(newTestUser(hashedTestPassword), nil)

		router := buildRouter(repo)

		rec := testutil.PerformRequest(t, router, http.MethodPost, loginEndpoint, requests.LogInRequest{
			Login:    "testuser",
			Password: testPassword,
		})

		assert.Equal(t, http.StatusOK, rec.Code)
		cookies := rec.Result().Cookies()
		assert.Len(t, cookies, 2)
		assert.NotNil(t, testutil.CookieByName(cookies, "accessToken"))
		assert.NotNil(t, testutil.CookieByName(cookies, "refreshToken"))
	})

	t.Run("FailureCases", func(t *testing.T) {
		type setupFn func(repo *mocks.MockRepository)

		cases := []struct {
			name          string
			setup         setupFn
			requestBody   any
			rawBody       []byte
			expectedCode  int
			expectedError string
		}{
			{
				name:          "InvalidJSON",
				rawBody:       []byte(`{"login": "testuser", "password": "oops"`),
				expectedCode:  http.StatusBadRequest,
				expectedError: "invalid request format",
			},
			{
				name: "UserNotFound",
				setup: func(repo *mocks.MockRepository) {
					repo.EXPECT().GetByLogin(gomock.Any(), "ghost").Return(nil, nil)
				},
				requestBody:   requests.LogInRequest{Login: "ghost", Password: testPassword},
				expectedCode:  http.StatusUnauthorized,
				expectedError: "invalid credentials",
			},
			{
				name: "WrongPassword",
				setup: func(repo *mocks.MockRepository) {
					repo.EXPECT().
						GetByLogin(gomock.Any(), "testuser").
						Return(newTestUser(mustHash("correctpassword")), nil)
				},
				requestBody:   requests.LogInRequest{Login: "testuser", Password: "wrongpassword"},
				expectedCode:  http.StatusUnauthorized,
				expectedError: "invalid credentials",
			},
			{
				name: "RepositoryInternalError",
				setup: func(repo *mocks.MockRepository) {
					repo.EXPECT().
						GetByLogin(gomock.Any(), "testuser").
						Return(nil, domain.ErrInternal)
				},
				requestBody:   requests.LogInRequest{Login: "testuser", Password: testPassword},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to log in",
			},
			{
				name: "RepositoryGenericError",
				setup: func(repo *mocks.MockRepository) {
					repo.EXPECT().
						GetByLogin(gomock.Any(), "testuser").
						Return(nil, errors.New("db down"))
				},
				requestBody:   requests.LogInRequest{Login: "testuser", Password: testPassword},
				expectedCode:  http.StatusInternalServerError,
				expectedError: "failed to log in",
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				repo := mocks.NewMockRepository(ctrl)
				if tc.setup != nil {
					tc.setup(repo)
				}

				router := buildRouter(repo)

				var rec *httptest.ResponseRecorder
				if tc.rawBody != nil {
					req := httptest.NewRequest(http.MethodPost, loginEndpoint, bytes.NewReader(tc.rawBody))
					req.Header.Set("Content-Type", "application/json")
					rec = httptest.NewRecorder()
					router.ServeHTTP(rec, req)
				} else {
					rec = testutil.PerformRequest(t, router, http.MethodPost, loginEndpoint, tc.requestBody)
				}

				require.Equal(t, tc.expectedCode, rec.Code)

				if tc.expectedError != "" {
					response := testutil.DecodeResponse(t, rec)
					assert.Equal(t, tc.expectedError, response["error"])
				}
			})
		}
	})
}

func newTestUser(hashedPassword string) *domain.User {
	return &domain.User{
		UUID:      uuid.New(),
		Login:     "testuser",
		Email:     "test@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func mustHash(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}
