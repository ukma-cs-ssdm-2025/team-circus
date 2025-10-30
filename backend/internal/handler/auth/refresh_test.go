package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"go.uber.org/zap"
)

type mockUserRepository struct {
	mock.Mock
}

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

	originalSecret := os.Getenv("SECRET_TOKEN")
	defer func() {
		if originalSecret == "" {
			os.Unsetenv("SECRET_TOKEN") //nolint:errcheck
		} else {
			os.Setenv("SECRET_TOKEN", originalSecret) //nolint:errcheck,gosec
		}
	}()

	os.Setenv("SECRET_TOKEN", "test-secret-key") //nolint:errcheck,gosec

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

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: tokenString, HttpOnly: true})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler(c)

		assert.Equal(t, http.StatusOK, w.Code)

		cookies := w.Result().Cookies()
		var newAccessCookie, newRefreshCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "accessToken" {
				newAccessCookie = cookie
			}
			if cookie.Name == "refreshToken" {
				newRefreshCookie = cookie
			}
		}

		assert.NotNil(t, newAccessCookie)
		assert.NotEmpty(t, newAccessCookie.Value)
		assert.True(t, newAccessCookie.HttpOnly)
		assert.Equal(t, "/", newAccessCookie.Path)
		assert.True(t, newAccessCookie.Secure)
		assert.Equal(t, http.SameSiteNoneMode, newAccessCookie.SameSite)

		assert.NotNil(t, newRefreshCookie)
		assert.NotEmpty(t, newRefreshCookie.Value)
		assert.True(t, newRefreshCookie.HttpOnly)
		assert.Equal(t, "/", newRefreshCookie.Path)
		assert.True(t, newRefreshCookie.Secure)
		assert.Equal(t, http.SameSiteNoneMode, newRefreshCookie.SameSite)

		repo.AssertExpectations(t)
	})

	t.Run("Missing refresh cookie", func(t *testing.T) {
		handler := auth.NewRefreshTokenHandler(new(mockUserRepository), zap.NewNop(), "test-secret-key", 10)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid token format", func(t *testing.T) {
		handler := auth.NewRefreshTokenHandler(new(mockUserRepository), zap.NewNop(), "test-secret-key", 10)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "invalid"})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Expired token", func(t *testing.T) {
		handler := auth.NewRefreshTokenHandler(new(mockUserRepository), zap.NewNop(), "test-secret-key", 10)

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		})
		tokenString, err := refreshToken.SignedString([]byte("test-secret-key"))
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: tokenString})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Missing secret token", func(t *testing.T) {
		os.Unsetenv("SECRET_TOKEN") //nolint:errcheck
		handler := auth.NewRefreshTokenHandler(new(mockUserRepository), zap.NewNop(), "test-secret-key", 10)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: "anything"})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		os.Setenv("SECRET_TOKEN", "test-secret-key") //nolint:errcheck,gosec
	})
}
