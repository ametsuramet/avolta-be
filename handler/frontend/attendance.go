package frontend

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TodayAttendanceHandler(c *gin.Context) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	getEmployee, _ := c.Get("employee")
	employee := getEmployee.(model.Employee)
	now := time.Now().In(loc)
	var data model.Attendance
	database.DB.Find(&data, "DATE(clock_in) = ? and employee_id = ?", now.Format("2006-01-02"), employee.ID)

	util.ResponsePaginatorSuccess(c, "Post Attandance Succeed", data, nil)
}
func ClockoutAttendanceHandler(c *gin.Context) {
	var data model.Attendance

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	var input model.AttendanceReq

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	data.ClockOut = input.ClockOut
	data.ClockOutLat = input.ClockOutLat
	data.ClockOutLng = input.ClockOutLng
	data.ClockOutNotes = util.SavedString(input.ClockOutNotes)
	data.ClockOutPicture = util.SavedString(input.ClockOutPicture)

	if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Checkout Attendance Succeed", nil, nil)

}
func PostAttendanceHandler(c *gin.Context) {
	getEmployee, _ := c.Get("employee")
	employee := getEmployee.(model.Employee)
	var input model.AttendanceReq

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Create(&model.Attendance{
		EmployeeID:     &employee.ID,
		ClockIn:        input.ClockIn,
		ClockInLat:     input.ClockInLat,
		ClockInLng:     input.ClockInLng,
		ClockInNotes:   util.SavedString(input.ClockInNotes),
		ClockInPicture: util.SavedString(input.ClockInPicture),
	}).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Post Attandance Succeed", nil, nil)
}
