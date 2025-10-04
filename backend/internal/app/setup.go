package app

import (
	"github.com/gin-gonic/gin"
)

func (a *App) setupRouter() *gin.Engine {
	router := gin.Default()

	_ = router.Group("/api")
	{

	}
	return router
}
