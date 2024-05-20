package router

import (
	"avolta/handler"
	"avolta/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3033", "https://avolta.web.app"},
		AllowMethods:  []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD"},
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Length", "Content-Type", "Access-Control-Allow-Origin", "API-KEY", "Currency-Code", "Cache-Control", "X-Requested-With", "Content-Disposition", "Content-Description"},
		ExposeHeaders: []string{"Content-Length", "Content-Disposition", "Content-Description"},
	}))

	r.Static("/assets", "./assets")

	v1 := r.Group("/api/v1")

	// v1.GET("/data", middleware.AuthMiddleware(), handler.GetData)

	admin := v1.Group("/admin")

	admin.Use()
	{
		admin.POST("/login", handler.Login)
	}

	admin.Use(middleware.AuthMiddleware())
	{

		category := admin.Group("/category")
		category.Use()
		{
			category.GET("", handler.CategoryGetAllHandler)
			category.GET("/:id", handler.CategoryGetOneHandler)
			category.POST("", handler.CategoryCreateHandler)
			category.PUT("/:id", handler.CategoryUpdateHandler)
			category.DELETE("/:id", handler.CategoryDeleteHandler)
		}

		my := admin.Group("/my")
		my.Use()
		{
			my.GET("", handler.MyController)
		}
		image := admin.Group("/image")
		image.Use()
		{
			image.GET("", handler.ImageGetAllHandler)
			image.GET("/:id", handler.ImageGetOneHandler)
			image.POST("", handler.ImageCreateHandler)
			image.PUT("/:id", handler.ImageUpdateHandler)
			image.DELETE("/:id", handler.ImageDeleteHandler)
		}

		file := admin.Group("/file")
		file.Use()
		{
			file.POST("/upload", handler.FileUploadHandler)
		}

		employee := admin.Group("/employee")
		employee.Use(middleware.AdminMiddleware())
		{
			employee.GET("", middleware.PermissionMiddleware("read_employee"), handler.EmployeeGetAllHandler)
			employee.POST("/import", middleware.PermissionMiddleware("import_employee"), handler.EmployeeImportHandler)
			employee.GET("/:id", middleware.PermissionMiddleware("read_employee"), handler.EmployeeGetOneHandler)
			employee.POST("", middleware.PermissionMiddleware("create_employee"), handler.EmployeeCreateHandler)
			employee.PUT("/:id", middleware.PermissionMiddleware("update_employee"), handler.EmployeeUpdateHandler)
			employee.DELETE("/:id", middleware.PermissionMiddleware("delete_employee"), handler.EmployeeDeleteHandler)
		}

		payRoll := admin.Group("/payRoll")
		payRoll.Use(middleware.AdminMiddleware())
		{
			payRoll.GET("", middleware.PermissionMiddleware("read_pay_roll"), handler.PayRollGetAllHandler)
			payRoll.GET("/:id", middleware.PermissionMiddleware("read_pay_roll"), handler.PayRollGetOneHandler)
			payRoll.POST("", middleware.PermissionMiddleware("create_pay_roll"), handler.PayRollCreateHandler)
			payRoll.PUT("/:id", middleware.PermissionMiddleware("update_pay_roll"), handler.PayRollUpdateHandler)
			payRoll.DELETE("/:id", middleware.PermissionMiddleware("delete_pay_roll"), handler.PayRollDeleteHandler)
		}

		leave := admin.Group("/leave")
		leave.Use(middleware.AdminMiddleware())
		{

			leave.GET("", middleware.PermissionMiddleware("read_leave"), handler.LeaveGetAllHandler)
			leave.GET("/:id", middleware.PermissionMiddleware("read_leave"), handler.LeaveGetOneHandler)
			leave.POST("", middleware.PermissionMiddleware("create_leave"), handler.LeaveCreateHandler)
			leave.PUT("/:id", middleware.PermissionMiddleware("update_leave"), handler.LeaveUpdateHandler)
			leave.DELETE("/:id", middleware.PermissionMiddleware("delete_leave"), handler.LeaveDeleteHandler)
		}

		attendance := admin.Group("/attendance")
		attendance.Use(middleware.AdminMiddleware())
		{
			attendance.GET("", middleware.PermissionMiddleware("read_attendance"), handler.AttendanceGetAllHandler)
			attendance.POST("/import", middleware.PermissionMiddleware("import_attendance"), handler.AttendanceImportHandler)
			attendance.GET("/import/:id", middleware.PermissionMiddleware("import_attendance"), handler.AttendanceImportDetailHandler)
			attendance.PUT("/import/:id/Reject", middleware.PermissionMiddleware("import_attendance_approval"), handler.AttendanceImportRejectHandler)
			attendance.PUT("/import/:id/Approve", middleware.PermissionMiddleware("import_attendance_approval"), handler.AttendanceImportApproveHandler)
			attendance.GET("/:id", middleware.PermissionMiddleware("read_attendance"), handler.AttendanceGetOneHandler)
			attendance.POST("", middleware.PermissionMiddleware("create_attendance"), handler.AttendanceCreateHandler)
			attendance.PUT("/:id", middleware.PermissionMiddleware("update_attendance"), handler.AttendanceUpdateHandler)
			attendance.DELETE("/:id", middleware.PermissionMiddleware("delete_attendance"), handler.AttendanceDeleteHandler)
		}

		incentive := admin.Group("/incentive")
		incentive.Use(middleware.AdminMiddleware())
		{
			incentive.GET("", handler.IncentiveGetAllHandler)
			incentive.GET("/:id", handler.IncentiveGetOneHandler)
			incentive.POST("", handler.IncentiveCreateHandler)
			incentive.PUT("/:id", handler.IncentiveUpdateHandler)
			incentive.DELETE("/:id", handler.IncentiveDeleteHandler)
		}

		account := admin.Group("/account")
		account.Use(middleware.AdminMiddleware())
		{
			account.GET("", handler.AccountGetAllHandler)
			account.GET("/:id", handler.AccountGetOneHandler)
			account.POST("", handler.AccountCreateHandler)
			account.PUT("/:id", handler.AccountUpdateHandler)
			account.DELETE("/:id", handler.AccountDeleteHandler)
		}

		permission := admin.Group("/permission")
		permission.Use(middleware.AdminMiddleware())
		{
			permission.GET("", handler.PermissionGetAllHandler)
			permission.GET("/:id", handler.PermissionGetOneHandler)
			// permission.POST("", handler.PermissionCreateHandler)
			// permission.PUT("/:id", handler.PermissionUpdateHandler)
			// permission.DELETE("/:id", handler.PermissionDeleteHandler)
		}

		transaction := admin.Group("/transaction")
		transaction.Use(middleware.AdminMiddleware())
		{
			transaction.GET("", handler.TransactionGetAllHandler)
			transaction.GET("/:id", handler.TransactionGetOneHandler)
			transaction.POST("", handler.TransactionCreateHandler)
			transaction.PUT("/:id", handler.TransactionUpdateHandler)
			transaction.DELETE("/:id", handler.TransactionDeleteHandler)
		}

		role := admin.Group("/role")
		role.Use(middleware.AdminMiddleware())
		{
			role.GET("", handler.RoleGetAllHandler)
			role.GET("/:id", handler.RoleGetOneHandler)
			role.POST("", handler.RoleCreateHandler)
			role.PUT("/:id", handler.RoleUpdateHandler)
			role.DELETE("/:id", handler.RoleDeleteHandler)
		}

		organization := admin.Group("/organization")
		organization.Use(middleware.AdminMiddleware())
		{
			organization.GET("", handler.OrganizationGetAllHandler)
			organization.GET("/:id", handler.OrganizationGetOneHandler)
			organization.POST("", handler.OrganizationCreateHandler)
			organization.PUT("/:id", handler.OrganizationUpdateHandler)
			organization.DELETE("/:id", handler.OrganizationDeleteHandler)
		}

		jobTitle := admin.Group("/jobTitle")
		jobTitle.Use(middleware.AdminMiddleware())
		{
			jobTitle.GET("", handler.JobTitleGetAllHandler)
			jobTitle.GET("/:id", handler.JobTitleGetOneHandler)
			jobTitle.POST("", handler.JobTitleCreateHandler)
			jobTitle.PUT("/:id", handler.JobTitleUpdateHandler)
			jobTitle.DELETE("/:id", handler.JobTitleDeleteHandler)
		}

		schedule := admin.Group("/schedule")
		schedule.Use(middleware.AdminMiddleware())
		{
			schedule.GET("", handler.ScheduleGetAllHandler)
			schedule.GET("/:id", handler.ScheduleGetOneHandler)
			schedule.PUT("/:id/AddEmployee", handler.ScheduleAddEmployeeHandler)
			schedule.DELETE("/:id/DeleteEmployee/:employeeId", handler.ScheduleDeleteEmployeeHandler)
			schedule.POST("", handler.ScheduleCreateHandler)
			schedule.PUT("/:id", handler.ScheduleUpdateHandler)
			schedule.DELETE("/:id", handler.ScheduleDeleteHandler)
		}

		leaveCategory := admin.Group("/leaveCategory")
		leaveCategory.Use(middleware.AdminMiddleware())
		{
			leaveCategory.GET("", middleware.PermissionMiddleware("read_leaved"), handler.LeaveCategoryGetAllHandler)
			leaveCategory.GET("/:id", middleware.PermissionMiddleware("read_leaved"), handler.LeaveCategoryGetOneHandler)
			leaveCategory.POST("", middleware.PermissionMiddleware("create_leaved"), handler.LeaveCategoryCreateHandler)
			leaveCategory.PUT("/:id", middleware.PermissionMiddleware("update_leaved"), handler.LeaveCategoryUpdateHandler)
			leaveCategory.DELETE("/:id", middleware.PermissionMiddleware("delete_leaved"), handler.LeaveCategoryDeleteHandler)
		}

		user := admin.Group("/user")
		user.Use(middleware.AdminMiddleware())
		{
			user.GET("", middleware.PermissionMiddleware("read_user"), handler.UserGetAllHandler)
			user.GET("/:id", middleware.PermissionMiddleware("read_user"), handler.UserGetOneHandler)
			user.POST("", middleware.PermissionMiddleware("create_user"), handler.UserCreateHandler)
			user.PUT("/:id", middleware.PermissionMiddleware("update_user"), handler.UserUpdateHandler)
			user.DELETE("/:id", middleware.PermissionMiddleware("delete_user"), handler.UserDeleteHandler)
		}

		// DONT REMOVE THIS LINE

	}
	return r
}
