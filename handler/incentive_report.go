package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IncentiveReportGetAllHandler(c *gin.Context) {
	var data []model.IncentiveReport

	preloads := []string{"User"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	// paginator.Joins = append(paginator.Joins, map[string]interface{}{
	// 	"JOIN shops ON shops.id = incentives.shop_id": nil,
	// })
	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }

	joinIncentive := false

	startDate, ok := c.GetQuery("start_date")
	if ok {
		paginator.WhereMoreEqual = append(paginator.WhereMoreEqual, map[string]interface{}{
			"inventive_reports.start_date": startDate,
		})

	}
	endDate, ok := c.GetQuery("end_date")
	if ok {
		paginator.WhereLessEqual = append(paginator.WhereLessEqual, map[string]interface{}{
			"inventive_reports.end_date": endDate,
		})

	}

	employeeId, ok := c.GetQuery("employee_id")
	if ok {
		joinIncentive = true

		paginator.Where = append(paginator.Where, map[string]interface{}{
			"incentives.employee_id": employeeId,
		})

	}
	employeeIds, ok := c.GetQuery("employee_ids")
	if ok {
		joinIncentive = true
		paginator.WhereIn = append(paginator.WhereIn,
			map[string][]string{
				"incentives.employee_id": strings.Split(employeeIds, ","),
			},
		)

	}
	// shopId, ok := c.GetQuery("shop_id")
	// if ok {
	// 	paginator.Where = append(paginator.Where, map[string]interface{}{
	// 		"incentives.shop_id": shopId,
	// 	})

	// }
	shopIds, ok := c.GetQuery("shop_ids")
	if ok {
		joinIncentive = true
		paginator.WhereIn = append(paginator.WhereIn,
			map[string][]string{
				"incentives.shop_id": strings.Split(shopIds, ","),
			},
		)

	}

	if joinIncentive {
		paginator.Joins = append(paginator.Joins, map[string]interface{}{
			"JOIN incentives ON incentives.incentive_report_id = incentive_reports.id": nil,
		})
		paginator.Joins = append(paginator.Joins, map[string]interface{}{
			"JOIN employees ON employees.id = incentives.employee_id": nil,
		})
	}

	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List IncentiveReport Retrived", dataRecords.Records, dataRecords)
}

func IncentiveReportGetOneHandler(c *gin.Context) {
	var data model.IncentiveReport

	id := c.Params.ByName("id")

	if err := database.DB.Preload("User").Preload("Shops").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveReport Retrived", data, nil)
}

func IncentiveReportCreateHandler(c *gin.Context) {

	var input model.IncentiveReportReq

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	getUser, _ := c.Get("user")
	user := getUser.(model.User)

	data := model.IncentiveReport{
		UserID:       user.ID,
		StartDate:    input.StartDate,
		EndDate:      input.EndDate,
		ReportNumber: input.ReportNumber,
		Description:  input.Description,
	}

	getCompany, _ := c.Get("company")
	company := getCompany.(model.Company)
	data.CompanyID = company.ID

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	shops := []model.Shop{}
	for _, v := range input.ShopIDs {
		shop := model.Shop{}
		database.DB.Find(&shop, "id = ?", v)
		shops = append(shops, shop)
	}

	database.DB.Model(&data).Association("Shops").Append(&shops)
	util.ResponseSuccess(c, "Data IncentiveReport Created", gin.H{"last_id": data.ID}, nil)
}

func IncentiveReportUpdateHandler(c *gin.Context) {
	var input, data model.IncentiveReport
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveReport Updated", nil, nil)
}

func IncentiveReportEditIncentiveHandler(c *gin.Context) {
	var data model.IncentiveReport
	id := c.Params.ByName("id")

	if err := database.DB.Preload("Shops").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := data.UpdateIncentive(c); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveReport Updated", nil, nil)
}

func IncentiveReportDeleteHandler(c *gin.Context) {
	var input, data model.IncentiveReport
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveReport Deleted", nil, nil)
}

func IncentiveReportAddEmployeeHandler(c *gin.Context) {
	var data model.IncentiveReport
	id := c.Params.ByName("id")

	if err := database.DB.Preload("Shops").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := data.AddEmployee(c); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveReport Updated", nil, nil)
}
