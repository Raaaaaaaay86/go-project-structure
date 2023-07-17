package http

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http/controller"
	"github.com/raaaaaaaay86/go-project-structure/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewHttpServer
//
//	@title						Video Service
//	@version					1.0
//	@description				This is a demo of Go project structure and others tech stack usage.
//	@contact.name				Ray Lin
//	@contact.email				raylincontact@icloud.com
//	@host						localhost:8080/video/api
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func NewHttpServer(
	authController controller.IAuthenticationController,
	videoController controller.IVideoController,
	commentController controller.ICommentController,
) *gin.Engine {
	engine := gin.Default()

	root := engine.Group("/video")
	{
		v1 := root.Group("/api/v1")
		{
			setAuthRouter(v1, authController)
			setVideoRouter(v1, videoController)
			setCommentRouter(v1, commentController)
		}
	}

	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}

func setAuthRouter(parent *gin.RouterGroup, auth controller.IAuthenticationController) {
	group := parent.Group("/auth")
	{
		group.POST("/register", auth.Register)
		group.POST("/login", auth.Login)
	}
}

func setVideoRouter(parent *gin.RouterGroup, video controller.IVideoController) {
	group := parent.Group("/video", middleware.Token)
	{
		group.POST("/upload", video.Upload)
		group.POST("/create", video.Create)
	}
}

func setCommentRouter(parent *gin.RouterGroup, comment controller.ICommentController) {
	group := parent.Group("/comment")
	{
		group.POST("/create", middleware.Token, comment.Create)
		group.GET("/find", comment.Find)
		group.DELETE("/delete", middleware.Token, comment.UserDelete)
		{
			adminGroup := group.Group("/admin", middleware.Token)
			adminGroup.DELETE("/delete", comment.ForceDelete)
		}
	}
}
