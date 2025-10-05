package app

import (
	"github.com/gin-gonic/gin"
	grouphandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	groupservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/group"
)

func (a *App) setupRouter() *gin.Engine {
	router := gin.Default()

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
