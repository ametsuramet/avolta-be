package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AttendanceGetAllHandler(c *gin.Context) {
	var data []model.Attendance
	preloads := []string{"Employee", "Employee.JobTitle"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	paginator.Paginate(&data)
	search, ok := c.GetQuery("search")

	paginator.Joins = append(paginator.Joins, map[string]interface{}{
		"LEFT JOIN employees ON employees.id = attendances.employee_id": nil,
	})

	if ok {
		paginator.Search = append(paginator.Search, map[string]interface{}{
			"employees.full_name":                search,
			"employees.employee_identity_number": search,
		})
	}
	jobTitleId, ok := c.GetQuery("job_title_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"employees.job_title_id": jobTitleId,
		})

	}
	gender, ok := c.GetQuery("gender")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"employees.gender": gender,
		})

	}
	startDate, ok := c.GetQuery("start_date")
	if ok {
		paginator.WhereMoreEqual = append(paginator.WhereMoreEqual, map[string]interface{}{
			"attendances.clock_in": startDate,
		})

	}
	endDate, ok := c.GetQuery("end_date")
	if ok {
		paginator.WhereLessEqual = append(paginator.WhereLessEqual, map[string]interface{}{
			"attendances.clock_in": endDate,
		})

	}
	employeeIds, ok := c.GetQuery("employee_ids")
	if ok {
		paginator.WhereIn = append(paginator.WhereIn,
			map[string][]string{
				"attendances.employee_id": strings.Split(employeeIds, ","),
			},
		)

	}

	// _, ok = c.GetQuery("download")
	// if ok {
	// 	paginator.Preloads = append(paginator.Preloads, "JobTitle")

	// }
	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponsePaginatorSuccess(c, "Data List Attandance Retrived", dataRecords.Records, dataRecords)
}

func AttendanceGetOneHandler(c *gin.Context) {
	var data model.Attendance

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Attendance Retrived", data, nil)
}

func AttendanceCreateHandler(c *gin.Context) {
	var data model.Attendance

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
	util.ResponseSuccess(c, "Data Attendance Created", gin.H{"last_id": data.ID}, nil)
}

func AttendanceUpdateHandler(c *gin.Context) {
	var input, data model.Attendance
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
	util.ResponseSuccess(c, "Data Attendance Updated", nil, nil)
}

func AttendanceDeleteHandler(c *gin.Context) {
	var input, data model.Attendance
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Attendance Deleted", nil, nil)
}
