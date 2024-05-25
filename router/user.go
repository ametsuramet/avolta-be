package router

import (
	"avolta/handler"
	"avolta/handler/frontend"
	"avolta/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(v1 *gin.RouterGroup) {
	user := v1.Group("/user")
	user.Use()
	{
		user.POST("/login", handler.LoginUser)
		user.POST("/upload", middleware.AuthMiddleware(), middleware.UserMiddleware(), handler.FileUploadHandler)
		my := user.Group("my")
		my.Use(middleware.AuthMiddleware(), middleware.UserMiddleware())
		{
			my.GET("", handler.MyController)
		}

		attendance := user.Group("/attendance")
		attendance.Use(middleware.AuthMiddleware(), middleware.UserMiddleware())
		{
			attendance.GET("/today", frontend.TodayAttendanceHandler)
			attendance.POST("", frontend.PostAttendanceHandler)
		}
	}
}
