package app

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	grouphandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	groupservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/group"
)

func (a *App) setupRouter() *gin.Engine {
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	groupRepo := grouprepo.NewGroupRepository(a.db)
	groupService := groupservice.NewGroupService(groupRepo)

	apiV1 := router.Group("/api/v1")
	{
		groups := apiV1.Group("/groups")
		{
			groups.POST("", grouphandler.NewCreateGroupHandler(groupService))
			groups.GET("/:uuid", grouphandler.NewGetGroupHandler(groupService))
			groups.GET("", grouphandler.NewGetAllGroupsHandler(groupService))
			groups.PUT("/:uuid", grouphandler.NewUpdateGroupHandler(groupService))
			groups.DELETE("/:uuid", grouphandler.NewDeleteGroupHandler(groupService))
		}
	}
	return router
}
