package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	"go.uber.org/zap"
)

func TestNewLogOutHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SuccessfulLogout", func(t *testing.T) {
		// Arrange
		handler := auth.NewLogOutHandler(zap.NewNop())

		req := httptest.NewRequest("POST", "/auth/logout", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Act
		handler(c)

		// Assert
		// Note: gin defaults to 200 when cookies are set, even though c.Status(204) is called
		assert.Equal(t, http.StatusOK, w.Code)

		// Check that cookies are expired
		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 2)

		var accessTokenCookie, refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			//nolint:staticcheck // single condition is clearer here
			if cookie.Name == "accessToken" {
				accessTokenCookie = cookie
			} else if cookie.Name == "refreshToken" {
				refreshTokenCookie = cookie
			}
		}

		assert.NotNil(t, accessTokenCookie)
		assert.NotNil(t, refreshTokenCookie)
		assert.Equal(t, "", accessTokenCookie.Value)
		assert.Equal(t, "", refreshTokenCookie.Value)
		assert.Equal(t, -1, accessTokenCookie.MaxAge)
		assert.Equal(t, -1, refreshTokenCookie.MaxAge)
	})
}
