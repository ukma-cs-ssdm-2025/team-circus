package app

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	documenthandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document"
	grouphandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group"
	reghandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg"
	userhandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user"
	documentrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/document"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	regrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/reg"
	userrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
	documentservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/document"
	groupservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/group"
	regservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/reg"
	userservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/user"
)

func (a *App) setupRouter() *gin.Engine {
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	groupRepo := grouprepo.NewGroupRepository(a.db)
	groupService := groupservice.NewGroupService(groupRepo)

	documentRepo := documentrepo.NewDocumentRepository(a.db)
	documentService := documentservice.NewDocumentService(documentRepo)

	userRepo := userrepo.NewUserRepository(a.db)
	userService := userservice.NewUserService(userRepo)

	regRepo := regrepo.NewRegRepository(a.db)
	regService := regservice.NewRegService(regRepo)

	apiV1 := router.Group("/api/v1")
	{
		groups := apiV1.Group("/groups")
		{
			groups.POST("", grouphandler.NewCreateGroupHandler(groupService))
			groups.GET("/:uuid", grouphandler.NewGetGroupHandler(groupService))
			groups.GET("", grouphandler.NewGetAllGroupsHandler(groupService))
			groups.PUT("/:uuid", grouphandler.NewUpdateGroupHandler(groupService))
			groups.DELETE("/:uuid", grouphandler.NewDeleteGroupHandler(groupService))
			groups.GET("/:uuid/documents", documenthandler.NewGetDocumentsByGroupHandler(documentService))
		}

		documents := apiV1.Group("/documents")
		{
			documents.POST("", documenthandler.NewCreateDocumentHandler(documentService))
			documents.GET("/:uuid", documenthandler.NewGetDocumentHandler(documentService))
			documents.GET("", documenthandler.NewGetAllDocumentsHandler(documentService))
			documents.PUT("/:uuid", documenthandler.NewUpdateDocumentHandler(documentService))
			documents.DELETE("/:uuid", documenthandler.NewDeleteDocumentHandler(documentService))
		}

		users := apiV1.Group("/users")
		{
			users.GET("/:uuid", userhandler.NewGetUserHandler(userService))
			users.GET("", userhandler.NewGetAllUsersHandler(userService))
			users.PUT("/:uuid", userhandler.NewUpdateUserHandler(userService))
			users.DELETE("/:uuid", userhandler.NewDeleteUserHandler(userService))
		}

		reg := apiV1.Group("/users")
		{
			reg.POST("", reghandler.NewRegHandler(regService))
		}
	}
	return router
}
