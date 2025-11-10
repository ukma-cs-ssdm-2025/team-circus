package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/testutil"
	"go.uber.org/zap"
)

type mockUserRepository struct {
	mock.Mock
}

const refreshEndpoint = "/auth/refresh"

func (m *mockUserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	panic("not implemented")
}

func (m *mockUserRepository) GetByUUID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck
}

func TestNewRefreshTokenHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful refresh", func(t *testing.T) {
		repo := new(mockUserRepository)
		handler := auth.NewRefreshTokenHandler(repo, zap.NewNop(), "test-secret-key", 10)

		userID := uuid.New()
		repo.On("GetByUUID", mock.Anything, userID).Return(&domain.User{UUID: userID}, nil)

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		})
		tokenString, err := refreshToken.SignedString([]byte("test-secret-key"))
		assert.NoError(t, err)

		c, w := testutil.NewRequestContext(t, http.MethodPost, refreshEndpoint)
		c.Request.AddCookie(&http.Cookie{Name: "refreshToken", Value: tokenString, HttpOnly: true})

		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)

		cookies := w.Result().Cookies()
		assert.NotNil(t, testutil.CookieByName(cookies, "accessToken"))
		assert.NotNil(t, testutil.CookieByName(cookies, "refreshToken"))

		repo.AssertExpectations(t)
	})

	t.Run("FailureCases", func(t *testing.T) {
		type refreshContextGen func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder)

		cases := []struct {
			name         string
			buildContext refreshContextGen
			expectedCode int
			setupHandler func() gin.HandlerFunc
		}{
			{
				name: "Missing refresh cookie",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					return testutil.NewRequestContext(t, http.MethodPost, refreshEndpoint)
				},
				expectedCode: http.StatusUnauthorized,
			},
			{
				name: "Invalid token format",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					c, w := testutil.NewRequestContext(t, http.MethodPost, refreshEndpoint)
					c.Request.AddCookie(&http.Cookie{Name: "refreshToken", Value: "invalid"})
					return c, w
				},
				expectedCode: http.StatusUnauthorized,
			},
			{
				name: "Expired token",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
						Subject:   uuid.NewString(),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
					})
					tokenString, err := token.SignedString([]byte("test-secret-key"))
					assert.NoError(t, err)

					c, w := testutil.NewRequestContext(t, http.MethodPost, refreshEndpoint)
					c.Request.AddCookie(&http.Cookie{Name: "refreshToken", Value: tokenString})
					return c, w
				},
				expectedCode: http.StatusUnauthorized,
			},
			{
				name: "Missing secret token",
				buildContext: func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
					c, w := testutil.NewRequestContext(t, http.MethodPost, refreshEndpoint)
					c.Request.AddCookie(&http.Cookie{Name: "refreshToken", Value: "anything"})
					return c, w
				},
				expectedCode: http.StatusInternalServerError,
				setupHandler: func() gin.HandlerFunc {
					return auth.NewRefreshTokenHandler(new(mockUserRepository), zap.NewNop(), "", 10)
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				handlerFactory := tc.setupHandler
				if handlerFactory == nil {
					handlerFactory = func() gin.HandlerFunc {
						return auth.NewRefreshTokenHandler(new(mockUserRepository), zap.NewNop(), "test-secret-key", 10)
					}
				}

				h := handlerFactory()
				c, w := tc.buildContext(t)

				h(c)

				assert.Equal(t, tc.expectedCode, w.Code)
			})
		}
	})
}
