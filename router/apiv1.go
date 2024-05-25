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
		AllowOrigins:  []string{"http://localhost:3033", "http://localhost:3034", "https://avolta.web.app", "https://avoltafe.web.app"},
		AllowMethods:  []string{"PUT", "PATCH", "GET", "POST", "DELETE", "HEAD"},
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Length", "Content-Type", "Access-Control-Allow-Origin", "API-KEY", "Currency-Code", "Cache-Control", "X-Requested-With", "Content-Disposition", "Content-Description"},
		ExposeHeaders: []string{"Content-Length", "Content-Disposition", "Content-Description"},
	}))

	r.Static("/assets", "./assets")

	v1 := r.Group("/api/v1")

	// v1.GET("/data", middleware.AuthMiddleware(), handler.GetData)

	UserRouter(v1)
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
			payRoll.PUT("/:id/Process", middleware.PermissionMiddleware("update_pay_roll"), handler.PayRollProcessHandler)
			payRoll.PUT("/:id/Payment", middleware.PermissionMiddleware("payment_pay_roll"), handler.PayRollPaymentHandler)
			payRoll.DELETE("/:id", middleware.PermissionMiddleware("delete_pay_roll"), handler.PayRollDeleteHandler)
		}
		payRollItem := admin.Group("/payRollItem")
		payRollItem.Use(middleware.AdminMiddleware())
		{
			// payRollItem.GET("", middleware.PermissionMiddleware("read_pay_roll"), handler.PayRollItemGetAllHandler)
			// payRollItem.GET("/:id", middleware.PermissionMiddleware("read_pay_roll"), handler.PayRollItemGetOneHandler)
			payRollItem.POST("", middleware.PermissionMiddleware("create_pay_roll"), handler.PayRollItemCreateHandler)
			payRollItem.PUT("/:id", middleware.PermissionMiddleware("update_pay_roll"), handler.PayRollItemUpdateHandler)
			payRollItem.DELETE("/:id", middleware.PermissionMiddleware("delete_pay_roll"), handler.PayRollItemDeleteHandler)
		}

		leave := admin.Group("/leave")
		leave.Use(middleware.AdminMiddleware())
		{

			leave.GET("", middleware.PermissionMiddleware("read_leave"), handler.LeaveGetAllHandler)
			leave.GET("/:id", middleware.PermissionMiddleware("read_leave"), handler.LeaveGetOneHandler)
			leave.POST("", middleware.PermissionMiddleware("create_leave"), handler.LeaveCreateHandler)
			leave.PUT("/:id", middleware.PermissionMiddleware("update_leave"), handler.LeaveUpdateHandler)
			leave.DELETE("/:id", middleware.PermissionMiddleware("delete_leave"), handler.LeaveDeleteHandler)
			leave.PUT("/:id/Approve", middleware.PermissionMiddleware("approval_leave"), handler.LeaveApproveHandler)
			leave.PUT("/:id/Reject", middleware.PermissionMiddleware("approval_leave"), handler.LeaveRejectHandler)
			leave.PUT("/:id/Review", middleware.PermissionMiddleware("update_leave"), handler.LeaveReviewHandler)
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
			leaveCategory.GET("", middleware.PermissionMiddleware("read_leave_category"), handler.LeaveCategoryGetAllHandler)
			leaveCategory.GET("/:id", middleware.PermissionMiddleware("read_leave_category"), handler.LeaveCategoryGetOneHandler)
			leaveCategory.POST("", middleware.PermissionMiddleware("create_leave_category"), handler.LeaveCategoryCreateHandler)
			leaveCategory.PUT("/:id", middleware.PermissionMiddleware("update_leave_category"), handler.LeaveCategoryUpdateHandler)
			leaveCategory.DELETE("/:id", middleware.PermissionMiddleware("delete_leave_category"), handler.LeaveCategoryDeleteHandler)
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

		company := admin.Group("/company")
		company.Use(middleware.AdminMiddleware())
		{
			company.GET("", middleware.PermissionMiddleware("menu_company"), handler.CompanyGetOneHandler)
			company.PUT("", middleware.PermissionMiddleware("menu_company"), handler.CompanyUpdateHandler)
		}

		setting := admin.Group("/setting")
		setting.Use(middleware.AdminMiddleware())
		{
			setting.GET("/autonumber", handler.SettingAutoNumberHandler)
			setting.GET("", handler.SettingGetOneHandler)
			setting.PUT("", middleware.PermissionMiddleware("menu_setting"), handler.SettingUpdateHandler)
		}

		reimbursement := admin.Group("/reimbursement")
		reimbursement.Use(middleware.AdminMiddleware())
		{
			reimbursement.GET("", middleware.PermissionMiddleware("read_reimbursement"), handler.ReimbursementGetAllHandler)
			reimbursement.GET("/:id", middleware.PermissionMiddleware("read_reimbursement"), handler.ReimbursementGetOneHandler)
			reimbursement.POST("", middleware.PermissionMiddleware("create_reimbursement"), handler.ReimbursementCreateHandler)
			reimbursement.PUT("/:id", middleware.PermissionMiddleware("update_reimbursement"), handler.ReimbursementUpdateHandler)
			reimbursement.PUT("/:id/Approval/:type", middleware.PermissionMiddleware("approvel_reimbursement"), handler.ReimbursementApprovalHandler)
			reimbursement.PUT("/:id/Payment", middleware.PermissionMiddleware("payment_reimbursement"), handler.ReimbursemenPaymentHandler)
			reimbursement.DELETE("/:id", middleware.PermissionMiddleware("delete_reimbursement"), handler.ReimbursementDeleteHandler)
		}

		reimbursementItem := admin.Group("/reimbursementItem")
		reimbursementItem.Use(middleware.AdminMiddleware())
		{
			reimbursementItem.GET("", handler.ReimbursementItemGetAllHandler)
			reimbursementItem.GET("/:id", handler.ReimbursementItemGetOneHandler)
			reimbursementItem.POST("", handler.ReimbursementItemCreateHandler)
			reimbursementItem.PUT("/:id", handler.ReimbursementItemUpdateHandler)
			reimbursementItem.DELETE("/:id", handler.ReimbursementItemDeleteHandler)
		}

		// payRollCost := admin.Group("/payRollCost")
		// payRollCost.Use()
		// {
		// 	payRollCost.GET("", handler.PayRollCostGetAllHandler)
		// 	payRollCost.GET("/:id", handler.PayRollCostGetOneHandler)
		// 	payRollCost.POST("", handler.PayRollCostCreateHandler)
		// 	payRollCost.PUT("/:id", handler.PayRollCostUpdateHandler)
		// 	payRollCost.DELETE("/:id", handler.PayRollCostDeleteHandler)
		// }

		product := admin.Group("/product")
		product.Use()
		{
			product.GET("", handler.ProductGetAllHandler)
			product.GET("/:id", handler.ProductGetOneHandler)
			product.POST("", handler.ProductCreateHandler)
			product.PUT("/:id", handler.ProductUpdateHandler)
			product.DELETE("/:id", handler.ProductDeleteHandler)
		}

		productCategory := admin.Group("/productCategory")
		productCategory.Use()
		{
			productCategory.GET("", handler.ProductCategoryGetAllHandler)
			productCategory.GET("/:id", handler.ProductCategoryGetOneHandler)
			productCategory.POST("", handler.ProductCategoryCreateHandler)
			productCategory.PUT("/:id", handler.ProductCategoryUpdateHandler)
			productCategory.DELETE("/:id", handler.ProductCategoryDeleteHandler)
		}

		shop := admin.Group("/shop")
		shop.Use()
		{
			shop.GET("", handler.ShopGetAllHandler)
			shop.GET("/:id", handler.ShopGetOneHandler)
			shop.POST("", handler.ShopCreateHandler)
			shop.PUT("/:id", handler.ShopUpdateHandler)
			shop.DELETE("/:id", handler.ShopDeleteHandler)
		}

		// DONT REMOVE THIS LINE

	}
	return r
}
