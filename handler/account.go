package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//
//
//
//
//

func AccountGetAllHandler(c *gin.Context) {
	var data []model.Account
	preloads := []string{}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }

	typeString, ok := c.GetQuery("type")
	if ok {
		paginator.WhereIn = append(paginator.WhereIn,
			map[string][]string{
				"type": strings.Split(typeString, ","),
			},
		)
	}

	isTaxString, ok := c.GetQuery("is_tax")
	if ok {
		paginator.Where = append(paginator.Where,
			map[string]interface{}{
				"is_tax": isTaxString,
			},
		)
	}

	cashflowGroupString, ok := c.GetQuery("cashflow_group")
	if ok {
		paginator.Where = append(paginator.Where,
			map[string]interface{}{
				"cashflow_group": cashflowGroupString,
			},
		)
	}

	cashflowSubGroupString, ok := c.GetQuery("cashflow_sub_group")
	if ok {
		paginator.WhereIn = append(paginator.WhereIn,
			map[string][]string{
				"cashflow_sub_group": strings.Split(cashflowSubGroupString, ","),
			},
		)
	}

	categoryString, ok := c.GetQuery("category")
	if ok {
		paginator.WhereIn = append(paginator.WhereIn,
			map[string][]string{
				"category": strings.Split(categoryString, ","),
			},
		)
	}

	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List Account Retrived", dataRecords.Records, dataRecords)
}

func AccountGetOneHandler(c *gin.Context) {
	var data model.Account

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Account Retrived", data, nil)
}

func AccountCreateHandler(c *gin.Context) {
	var data model.Account

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
	util.ResponseSuccess(c, "Data Account Created", gin.H{"last_id": data.ID}, nil)
}

func AccountUpdateHandler(c *gin.Context) {
	var input, data model.Account
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
	util.ResponseSuccess(c, "Data Account Updated", nil, nil)
}

func AccountDeleteHandler(c *gin.Context) {
	var input, data model.Account
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Account Deleted", nil, nil)
}
