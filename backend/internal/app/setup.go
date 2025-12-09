package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	authhandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/auth"
	documenthandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/document"
	grouphandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group"
	memberhandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/member"
	reghandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg"
	userhandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user"
	websockethandler "github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/websocket"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/middleware"
	collabrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo"
	documentrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/document"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	memberrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/member"
	regrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/reg"
	userrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
	documentservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/document"
	groupservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/group"
	memberservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/member"
	regservice "github.com/ukma-cs-ssdm-2025/team-circus/internal/service/reg"
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

	groupRepo := grouprepo.NewGroupRepository(a.DB)
	memberRepo := memberrepo.NewMemberRepository(a.DB)
	groupService := groupservice.NewGroupService(groupRepo, memberRepo)

	documentRepo := documentrepo.NewDocumentRepository(a.DB)
	documentService := documentservice.NewDocumentService(documentRepo, memberRepo)
	documentPersistence := collabrepo.NewDocumentPersistence(a.DB)

	userRepo := userrepo.NewUserRepository(a.DB)
	userService := userservice.NewUserService(userRepo, a.cfg.HashingCost)

	regRepo := regrepo.NewRegRepository(a.DB)
	regService := regservice.NewRegService(regRepo, a.cfg.HashingCost)

	memberService := memberservice.NewMemberService(memberRepo, groupRepo, userRepo)

	apiV1 := router.Group("/api/v1")

	wsHubManager := websockethandler.NewHubManager(a.l, documentPersistence)

	public := apiV1.Group("")
	{
		public.POST("/signup", reghandler.NewRegHandler(regService, a.l))
		public.POST("/auth/login", authhandler.NewLogInHandler(userRepo, a.l,
			a.cfg.SecretToken, a.cfg.AccessDuration, a.cfg.RefreshDuration))
		public.POST("/auth/refresh", authhandler.NewRefreshTokenHandler(userRepo, a.l, a.cfg.SecretToken, a.cfg.AccessDuration))
	}

	protected := apiV1.Group("")
	protected.Use(middleware.AuthMiddleware(userRepo, a.cfg.SecretToken))
	{
		protected.POST("/auth/logout", authhandler.NewLogOutHandler(a.l))

		groups := protected.Group("/groups")
		{
			groups.POST("", grouphandler.NewCreateGroupHandler(groupService, a.l))
			groups.GET("/:uuid", grouphandler.NewGetGroupHandler(groupService, a.l))
			groups.GET("", grouphandler.NewGetAllGroupsHandler(groupService, a.l))
			groups.PUT("/:uuid", grouphandler.NewUpdateGroupHandler(groupService, a.l))
			groups.DELETE("/:uuid", grouphandler.NewDeleteGroupHandler(groupService, a.l))

			members := groups.Group("/:uuid/members")
			{
				members.GET("", memberhandler.NewGetAllMembersHandler(memberService, a.l))
				members.POST("", memberhandler.NewCreateMemberHandler(memberService, a.l))
				members.PUT("/:user_uuid", memberhandler.NewUpdateMemberHandler(memberService, a.l))
				members.DELETE("/:user_uuid", memberhandler.NewDeleteMemberHandler(memberService, a.l))
			}
		}

		documents := protected.Group("/documents")
		{
			documents.POST("", documenthandler.NewCreateDocumentHandler(documentService, a.l))
			documents.GET("/:uuid", documenthandler.NewGetDocumentHandler(documentService, a.l))
			documents.GET("", documenthandler.NewGetAllDocumentsHandler(documentService, a.l))
			documents.PUT("/:uuid", documenthandler.NewUpdateDocumentHandler(documentService, a.l))
			documents.DELETE("/:uuid", documenthandler.NewDeleteDocumentHandler(documentService, a.l))
		}

		users := protected.Group("/users")
		{
			users.GET("/:uuid", userhandler.NewGetUserHandler(userService, a.l))
			users.GET("", userhandler.NewGetAllUsersHandler(userService, a.l))
			users.PUT("/:uuid", userhandler.NewUpdateUserHandler(userService, a.l))
			users.DELETE("/:uuid", userhandler.NewDeleteUserHandler(userService, a.l))
		}
	}

	wsGroup := router.Group("/ws")
	wsGroup.Use(middleware.AuthMiddleware(userRepo, a.cfg.SecretToken))
	{
		wsGroup.GET("/documents/:uuid", websockethandler.NewWebSocketHandler(
			documentService,
			wsHubManager,
			a.l,
		))
	}
	return router
}
