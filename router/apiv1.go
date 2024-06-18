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

	admin.Use(middleware.GeneralMiddleware())
	{
		admin.POST("/login", handler.Login)
		admin.POST("/register", handler.Register)
		admin.POST("/verification/:token", handler.Verification)
	}

	admin.Use(middleware.AuthMiddleware())
	{
		admin.POST("/create/company", handler.CreateCompany)
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
			payRoll.PUT("/:id/Finish", middleware.PermissionMiddleware("update_pay_roll"), handler.FinishProcessHandler)
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
			attendance.GET("/summary/:employeeId", middleware.PermissionMiddleware("read_attendance"), handler.SummaryGetAllHandler)
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
			role.GET("", middleware.PermissionMiddleware("read_role"), handler.RoleGetAllHandler)
			role.GET("/:id", middleware.PermissionMiddleware("read_role"), handler.RoleGetOneHandler)
			role.POST("", middleware.PermissionMiddleware("create_role"), handler.RoleCreateHandler)
			role.PUT("/:id", middleware.PermissionMiddleware("update_role"), handler.RoleUpdateHandler)
			role.DELETE("/:id", middleware.PermissionMiddleware("delete_role"), handler.RoleDeleteHandler)
		}

		organization := admin.Group("/organization")
		organization.Use(middleware.AdminMiddleware())
		{
			organization.GET("", middleware.PermissionMiddleware("read_organization"), handler.OrganizationGetAllHandler)
			organization.GET("/:id", middleware.PermissionMiddleware("read_organization"), handler.OrganizationGetOneHandler)
			organization.POST("", middleware.PermissionMiddleware("create_organization"), handler.OrganizationCreateHandler)
			organization.PUT("/:id", middleware.PermissionMiddleware("update_organization"), handler.OrganizationUpdateHandler)
			organization.DELETE("/:id", middleware.PermissionMiddleware("delete_organization"), handler.OrganizationDeleteHandler)
		}

		jobTitle := admin.Group("/jobTitle")
		jobTitle.Use(middleware.AdminMiddleware())
		{
			jobTitle.GET("", middleware.PermissionMiddleware("read_job_title"), handler.JobTitleGetAllHandler)
			jobTitle.GET("/:id", middleware.PermissionMiddleware("read_job_title"), handler.JobTitleGetOneHandler)
			jobTitle.POST("", middleware.PermissionMiddleware("create_job_title"), handler.JobTitleCreateHandler)
			jobTitle.PUT("/:id", middleware.PermissionMiddleware("update_job_title"), handler.JobTitleUpdateHandler)
			jobTitle.DELETE("/:id", middleware.PermissionMiddleware("delete_job_title"), handler.JobTitleDeleteHandler)
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
			setting.GET("/incentive/autonumber", handler.SettingIncentiveAutoNumberHandler)
			setting.GET("/pay_roll_report/autonumber", handler.SettingPayrollReportAutoNumberHandler)
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
		product.Use(middleware.AdminMiddleware())
		{
			product.GET("", middleware.PermissionMiddleware("read_product"), handler.ProductGetAllHandler)
			product.GET("/:id", middleware.PermissionMiddleware("read_product"), handler.ProductGetOneHandler)
			product.POST("", middleware.PermissionMiddleware("create_product"), handler.ProductCreateHandler)
			product.PUT("/:id", middleware.PermissionMiddleware("update_product"), handler.ProductUpdateHandler)
			product.DELETE("/:id", middleware.PermissionMiddleware("delete_product"), handler.ProductDeleteHandler)
			product.POST("/import", middleware.PermissionMiddleware("import_product"), handler.ProductImportHandler)
		}

		productCategory := admin.Group("/productCategory")
		productCategory.Use(middleware.AdminMiddleware())
		{
			productCategory.GET("", middleware.PermissionMiddleware("read_product_category"), handler.ProductCategoryGetAllHandler)
			productCategory.GET("/:id", middleware.PermissionMiddleware("read_product_category"), handler.ProductCategoryGetOneHandler)
			productCategory.POST("", middleware.PermissionMiddleware("create_product_category"), handler.ProductCategoryCreateHandler)
			productCategory.PUT("/:id", middleware.PermissionMiddleware("update_product_category"), handler.ProductCategoryUpdateHandler)
			productCategory.DELETE("/:id", middleware.PermissionMiddleware("delete_product_category"), handler.ProductCategoryDeleteHandler)
		}

		shop := admin.Group("/shop")
		shop.Use(middleware.AdminMiddleware())
		{
			shop.GET("", middleware.PermissionMiddleware("read_shop"), handler.ShopGetAllHandler)
			shop.GET("/:id", middleware.PermissionMiddleware("read_shop"), handler.ShopGetOneHandler)
			shop.POST("", middleware.PermissionMiddleware("create_shop"), handler.ShopCreateHandler)
			shop.PUT("/:id", middleware.PermissionMiddleware("update_shop"), handler.ShopUpdateHandler)
			shop.DELETE("/:id", middleware.PermissionMiddleware("delete_shop"), handler.ShopDeleteHandler)
		}

		sale := admin.Group("/sale")
		sale.Use(middleware.AdminMiddleware())
		{
			sale.GET("", middleware.PermissionMiddleware("read_sale"), handler.SaleGetAllHandler)
			sale.GET("/:id", middleware.PermissionMiddleware("read_sale"), handler.SaleGetOneHandler)
			sale.POST("", middleware.PermissionMiddleware("create_sale"), handler.SaleCreateHandler)
			sale.PUT("/:id", middleware.PermissionMiddleware("update_sale"), handler.SaleUpdateHandler)
			sale.DELETE("/:id", middleware.PermissionMiddleware("delete_sale"), handler.SaleDeleteHandler)
			sale.POST("/import", middleware.PermissionMiddleware("import_sale"), handler.SaleImportHandler)
		}

		saleReceipt := admin.Group("/saleReceipt")
		saleReceipt.Use(middleware.AdminMiddleware())
		{
			saleReceipt.GET("", handler.SaleReceiptGetAllHandler)
			saleReceipt.GET("/:id", handler.SaleReceiptGetOneHandler)
			saleReceipt.POST("", handler.SaleReceiptCreateHandler)
			saleReceipt.PUT("/:id", handler.SaleReceiptUpdateHandler)
			saleReceipt.DELETE("/:id", handler.SaleReceiptDeleteHandler)
		}

		incentiveSetting := admin.Group("/incentiveSetting")
		incentiveSetting.Use(middleware.AdminMiddleware())
		{
			incentiveSetting.GET("", middleware.PermissionMiddleware("read_incentive_setting"), handler.IncentiveSettingGetAllHandler)
			incentiveSetting.GET("/:id", middleware.PermissionMiddleware("read_incentive_setting"), handler.IncentiveSettingGetOneHandler)
			incentiveSetting.POST("", middleware.PermissionMiddleware("create_incentive_setting"), handler.IncentiveSettingCreateHandler)
			incentiveSetting.PUT("/:id", middleware.PermissionMiddleware("update_incentive_setting"), handler.IncentiveSettingUpdateHandler)
			incentiveSetting.DELETE("/:id", middleware.PermissionMiddleware("delete_incentive_setting"), handler.IncentiveSettingDeleteHandler)
			incentiveSetting.POST("/import", middleware.PermissionMiddleware("import_incentive"), handler.IncentiveImportHandler)
		}

		incentiveShop := admin.Group("/incentiveShop")
		incentiveShop.Use(middleware.AdminMiddleware())
		{
			incentiveShop.GET("", middleware.PermissionMiddleware("read_incentive_shop"), handler.IncentiveShopGetAllHandler)
			incentiveShop.GET("/:id", middleware.PermissionMiddleware("read_incentive_shop"), handler.IncentiveShopGetOneHandler)
			// incentiveShop.POST("", middleware.PermissionMiddleware("create_incentive_shop"),  handler.IncentiveShopCreateHandler)
			// incentiveShop.PUT("/:id", middleware.PermissionMiddleware("update_incentive_shop"),  handler.IncentiveShopUpdateHandler)
			incentiveShop.DELETE("/:id", middleware.PermissionMiddleware("delete_incentive_shop"), handler.IncentiveShopDeleteHandler)
		}

		incentiveReport := admin.Group("/incentiveReport")
		incentiveReport.Use(middleware.AdminMiddleware())
		{
			incentiveReport.GET("", middleware.PermissionMiddleware("read_incentive_report"), handler.IncentiveReportGetAllHandler)
			incentiveReport.GET("/:id", middleware.PermissionMiddleware("read_incentive_report"), handler.IncentiveReportGetOneHandler)
			incentiveReport.POST("", middleware.PermissionMiddleware("create_incentive_report"), handler.IncentiveReportCreateHandler)
			incentiveReport.PUT("/:id", middleware.PermissionMiddleware("update_incentive_report"), handler.IncentiveReportUpdateHandler)
			incentiveReport.DELETE("/:id", middleware.PermissionMiddleware("delete_incentive_report"), handler.IncentiveReportDeleteHandler)
			incentiveReport.PUT("/:id/EditIncentive/:incentiveId", handler.IncentiveReportEditIncentiveHandler)
			incentiveReport.PUT("/:id/AddEmployee", handler.IncentiveReportAddEmployeeHandler)
		}

		bank := admin.Group("/bank")
		bank.Use(middleware.AdminMiddleware())
		{
			bank.GET("", middleware.PermissionMiddleware("read_bank"), handler.BankGetAllHandler)
			bank.GET("/:id", middleware.PermissionMiddleware("read_bank"), handler.BankGetOneHandler)
			bank.POST("", middleware.PermissionMiddleware("create_bank"), handler.BankCreateHandler)
			bank.PUT("/:id", middleware.PermissionMiddleware("update_bank"), handler.BankUpdateHandler)
			bank.DELETE("/:id", middleware.PermissionMiddleware("delete_bank"), handler.BankDeleteHandler)
		}

		announcement := admin.Group("/announcement")
		announcement.Use(middleware.AdminMiddleware())
		{
			announcement.GET("", middleware.PermissionMiddleware("read_announcement"), handler.AnnouncementGetAllHandler)
			announcement.GET("/:id", middleware.PermissionMiddleware("read_announcement"), handler.AnnouncementGetOneHandler)
			announcement.POST("", middleware.PermissionMiddleware("create_announcement"), handler.AnnouncementCreateHandler)
			announcement.PUT("/:id", middleware.PermissionMiddleware("update_announcement"), handler.AnnouncementUpdateHandler)
			announcement.DELETE("/:id", middleware.PermissionMiddleware("delete_announcement"), handler.AnnouncementDeleteHandler)
		}

		payRollReport := admin.Group("/payRollReport")
		payRollReport.Use(middleware.AdminMiddleware())
		{
			payRollReport.GET("", middleware.PermissionMiddleware("read_pay_roll_report"), handler.PayRollReportGetAllHandler)
			payRollReport.GET("/:id", middleware.PermissionMiddleware("read_pay_roll_report"), handler.PayRollReportGetOneHandler)
			payRollReport.POST("", middleware.PermissionMiddleware("create_pay_roll_report"), handler.PayRollReportCreateHandler)
			payRollReport.PUT("/:id", middleware.PermissionMiddleware("update_pay_roll_report"), handler.PayRollReportUpdateHandler)
			payRollReport.DELETE("/:id", middleware.PermissionMiddleware("delete_pay_roll_report"), handler.PayRollReportDeleteHandler)
			payRollReport.PUT("/:id/AddItem/:payrollID", middleware.PermissionMiddleware("update_pay_roll_report"), handler.PayRollReportAddItemHandler)
			payRollReport.GET("/:id/PayRolllBankDownLoad", middleware.PermissionMiddleware("update_pay_roll_report"), handler.PayRollReporDownloadPayRollBankHandler)
			payRollReport.DELETE("/:id/DeleteItem/:payrollID", middleware.PermissionMiddleware("update_pay_roll_report"), handler.PayRollReportDeleteItemHandler)
		}

		// payRollReportItem := admin.Group("/payRollReportItem")
		// payRollReportItem.Use()
		// {
		// 	payRollReportItem.GET("", handler.PayRollReportItemGetAllHandler)
		// 	payRollReportItem.GET("/:id", handler.PayRollReportItemGetOneHandler)
		// 	payRollReportItem.POST("", handler.PayRollReportItemCreateHandler)
		// 	payRollReportItem.PUT("/:id", handler.PayRollReportItemUpdateHandler)
		// 	payRollReportItem.DELETE("/:id", handler.PayRollReportItemDeleteHandler)
		// }

		// DONT REMOVE THIS LINE

	}
	return r
}
