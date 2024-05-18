package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ScheduleGetAllHandler(c *gin.Context) {
	var data []model.Schedule
	preloads := []string{"Employees"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	paginator.Paginate(&data)
	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }
	dateRange, ok := c.GetQuery("date_range")
	if ok {
		split := strings.Split(dateRange, ",")

		paginator.WhereQuery = append(paginator.WhereQuery, map[string][]interface{}{
			"(start_date >= ? and start_date <= ?) OR (schedule_type = 'WEEKLY')": {split[0], split[1]},
		})
	}
	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List Schedule Retrived", dataRecords.Records, dataRecords)
}

func ScheduleGetOneHandler(c *gin.Context) {
	var data model.Schedule

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Schedule Retrived", data, nil)
}

func ScheduleCreateHandler(c *gin.Context) {
	var data model.Schedule

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// getUser, _ := c.Get("user")
	// user := getUser.(model.User)

	// data.AuthorID = user.ID

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	employees := []model.Employee{}
	for _, v := range data.EmployeeIDs {
		employee := model.Employee{}
		database.DB.Find(&employee, "id = ?", v)
		employees = append(employees, employee)
	}
	database.DB.Model(&data).Association("Employees").Append(employees)
	util.ResponseSuccess(c, "Data Schedule Created", gin.H{"last_id": data.ID}, nil)
}

func ScheduleUpdateHandler(c *gin.Context) {
	var input, data model.Schedule
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
	util.ResponseSuccess(c, "Data Schedule Updated", nil, nil)
}

func ScheduleDeleteHandler(c *gin.Context) {
	var input, data model.Schedule
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Schedule Deleted", nil, nil)
}
