package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	documenthandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document"
	grouphandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group"
	userhandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user"
	documentrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/document"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	userrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
	documentservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/document"
	groupservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/group"
	userservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/user"
)

func (a *App) setupRouter() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     a.cfg.CORS.AllowOrigins,
		AllowMethods:     a.cfg.CORS.AllowMethods,
		AllowHeaders:     a.cfg.CORS.AllowHeaders,
		ExposeHeaders:    a.cfg.CORS.ExposeHeaders,
		AllowCredentials: a.cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(a.cfg.CORS.MaxAge) * time.Second,
	}))

	// Swagger documentation
	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	groupRepo := grouprepo.NewGroupRepository(a.db)
	groupService := groupservice.NewGroupService(groupRepo)

	documentRepo := documentrepo.NewDocumentRepository(a.db)
	documentService := documentservice.NewDocumentService(documentRepo)

	userRepo := userrepo.NewUserRepository(a.db)
	userService := userservice.NewUserService(userRepo)

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
			users.POST("", userhandler.NewCreateUserHandler(userService))
			users.GET("/:uuid", userhandler.NewGetUserHandler(userService))
			users.GET("", userhandler.NewGetAllUsersHandler(userService))
			users.PUT("/:uuid", userhandler.NewUpdateUserHandler(userService))
			users.DELETE("/:uuid", userhandler.NewDeleteUserHandler(userService))
		}
	}
	return router
}
