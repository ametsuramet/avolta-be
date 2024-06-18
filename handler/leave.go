package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LeaveGetAllHandler(c *gin.Context) {
	var data []model.Leave
	preloads := []string{"LeaveCategory"}
	paginator := util.NewPaginator(c)

	_, ok := c.GetQuery("show_employee")
	if ok {
		preloads = append(preloads, "Employee")
	}
	paginator.Preloads = preloads
	paginator.Joins = append(paginator.Joins, map[string]interface{}{
		"JOIN leave_categories ON leave_categories.id = leaves.leave_category_id": nil,
	})

	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }

	startDate, ok := c.GetQuery("start_date")
	if ok {
		paginator.WhereMoreEqual = append(paginator.WhereMoreEqual, map[string]interface{}{
			"leaves.start_date": startDate,
		})

	}
	endDate, ok := c.GetQuery("end_date")
	if ok {
		paginator.WhereLessEqual = append(paginator.WhereLessEqual, map[string]interface{}{
			"leaves.start_date": endDate,
		})

	}

	employeeId, ok := c.GetQuery("employee_id")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"leaves.employee_id": employeeId,
		})

	}
	status, ok := c.GetQuery("status")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"leaves.status": status,
		})

	}

	absent, ok := c.GetQuery("absent")
	if ok {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"leave_categories.absent": absent,
		})

	} else {
		paginator.Where = append(paginator.Where, map[string]interface{}{
			"leave_categories.absent": false,
		})
	}
	orderBy, ok := c.GetQuery("order_by")
	if ok {
		paginator.OrderBy = []string{orderBy}
	}

	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List Leave Retrived", dataRecords.Records, dataRecords)
}

func LeaveGetOneHandler(c *gin.Context) {
	var data model.Leave

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Leave Retrived", data, nil)
}

func LeaveCreateHandler(c *gin.Context) {
	var data model.Leave

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// getUser, _ := c.Get("user")
	// user := getUser.(model.User)

	// data.AuthorID = user.ID

	getCompany, _ := c.Get("company")
	company := getCompany.(model.Company)
	data.CompanyID = company.ID

	leaveCat := model.LeaveCategory{}
	database.DB.Find(&leaveCat, "id = ?", data.LeaveCategoryID)
	if leaveCat.Absent {
		data.Status = "APPROVED"
		data.RequestType = "FULL_DAY"
	}

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Leave Created", gin.H{"last_id": data.ID}, nil)
}

func LeaveUpdateHandler(c *gin.Context) {
	var input, data model.Leave
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
	util.ResponseSuccess(c, "Data Leave Updated", nil, nil)
}

func LeaveDeleteHandler(c *gin.Context) {
	var input, data model.Leave
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Leave Deleted", nil, nil)
}

func LeaveApproveHandler(c *gin.Context) {
	if err := approval(c, "APPROVED"); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Leave Approved", nil, nil)
}

func LeaveRejectHandler(c *gin.Context) {
	if err := approval(c, "REJECTED"); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Leave Rejected", nil, nil)
}
func LeaveReviewHandler(c *gin.Context) {
	if err := approval(c, "REVIEWED"); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Leave Rejected", nil, nil)
}

func approval(c *gin.Context, status string) error {
	var data model.Leave
	var input = struct {
		Remarks string `json:"remarks"`
	}{}
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		return err
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		return err
	}

	getUser, _ := c.Get("user")
	user := getUser.(model.User)
	data.ApproverID = &user.ID

	data.Status = status
	data.Remarks += "- **[" + user.FullName + "]** \n*" + time.Now().Format("02/01/2006 15:04") + "*\n\n" + input.Remarks + "\n\n"
	if err := database.DB.Model(&data).Updates(&data).Error; err != nil {
		return err
	}

	return nil
}
