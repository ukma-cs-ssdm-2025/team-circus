package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewLogOutHandler logs the user out by expiring auth cookies.
// @Summary User logout
// @Description Expires the JWT access/refresh cookies.
// @Tags auth
// @Success 204 "Logged out successfully"
// @Router /auth/logout [post]
func NewLogOutHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("accessToken", "", -1, "/", "", true, true)
		c.SetCookie("refreshToken", "", -1, "/", "", true, true)

		c.Status(http.StatusNoContent)
	}
}
