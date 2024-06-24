package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserGetAllHandler(c *gin.Context) {
	getCompany, _ := c.Get("company")
	company := getCompany.(model.Company)
	var data []model.User
	preloads := []string{"Employee", "Roles"}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads
	paginator.Joins = append(paginator.Joins, map[string]interface{}{
		"JOIN user_companies ON user_companies.user_id = users.id": nil,
	})
	paginator.Joins = append(paginator.Joins, map[string]interface{}{
		"JOIN companies ON companies.id = user_companies.company_id": nil,
	})

	paginator.Where = append(paginator.Where, map[string]interface{}{
		"companies.id": company.ID,
	})

	search, ok := c.GetQuery("search")
	if ok {
		paginator.Search = append(paginator.Search, map[string]interface{}{
			"full_name": search,
		})
	}
	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// fmt.Println("company", company)
	for i := range data {
		data[i].CompanyID = company.ID

	}

	util.ResponsePaginatorSuccess(c, "Data List User Retrived", dataRecords.Records, dataRecords)
}

func UserGetOneHandler(c *gin.Context) {
	var data model.User

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data User Retrived", data, nil)
}

func UserCreateHandler(c *gin.Context) {
	var data model.User

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
	util.ResponseSuccess(c, "Data User Created", gin.H{"last_id": data.ID}, nil)
}

func UserUpdateHandler(c *gin.Context) {
	var input, data model.User
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Preload("Employee").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.EmployeeID != "" && data.EmployeeID != input.EmployeeID {
		employee := model.Employee{}
		database.DB.Find(&employee, "id = ?", input.EmployeeID)
		database.DB.Model(&data).Association("Employee").Replace(&employee)
	}
	util.ResponseSuccess(c, "Data User Updated", nil, nil)
}

func UserDeleteHandler(c *gin.Context) {
	var input, data model.User
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data User Deleted", nil, nil)
}
