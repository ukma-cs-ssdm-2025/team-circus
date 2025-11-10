package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/mocks"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/testutil"
	"go.uber.org/zap"
)

const refreshEndpoint = "/auth/refresh"

func TestNewRefreshTokenHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful refresh", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mocks.NewMockRepository(ctrl)
		handler := auth.NewRefreshTokenHandler(repo, zap.NewNop(), "test-secret-key", 10)

		userID := uuid.New()
		repo.EXPECT().
			GetByUUID(gomock.Any(), userID).
			Return(&domain.User{UUID: userID}, nil)

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

	})

	t.Run("FailureCases", func(t *testing.T) {
		type refreshContextGen func(t *testing.T) (*gin.Context, *httptest.ResponseRecorder)

		cases := []struct {
			name         string
			buildContext refreshContextGen
			expectedCode int
			setupHandler func(ctrl *gomock.Controller) gin.HandlerFunc
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
				setupHandler: func(ctrl *gomock.Controller) gin.HandlerFunc {
					return auth.NewRefreshTokenHandler(mocks.NewMockRepository(ctrl), zap.NewNop(), "", 10)
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				handlerFactory := tc.setupHandler
				if handlerFactory == nil {
					handlerFactory = func(ctrl *gomock.Controller) gin.HandlerFunc {
						return auth.NewRefreshTokenHandler(mocks.NewMockRepository(ctrl), zap.NewNop(), "test-secret-key", 10)
					}
				}

				h := handlerFactory(ctrl)
				c, w := tc.buildContext(t)

				h(c)

				assert.Equal(t, tc.expectedCode, w.Code)
			})
		}
	})
}
