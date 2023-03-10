package routes

import (
	"oprec/go-blog/controller"
	"oprec/go-blog/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController controller.UserController) {
	userRegist := router.Group("/masuk")
	{
		userRegist.POST("/register", userController.Register)
		userRegist.POST("/login", userController.Login)
	}

	userGuest := router.Group("public")
	{
		userGuest.GET("/blogs", userController.GetAllBlog)
		userGuest.GET("/blogs/:id", userController.BlogDetails)
	}

	userRoutes := router.Group("/secured").Use(middleware.Authenticate())
	{
		userRoutes.GET("/myblogs", userController.GetBlogByID)
		userRoutes.POST("/upload", userController.Upload)
		userRoutes.PUT("/update/name", userController.Update)
		userRoutes.POST("/comment/:id", userController.Comment)
		userRoutes.POST("/like/blog/:id", userController.Like)
		userRoutes.POST("/delete", userController.Delete)
	}

}
