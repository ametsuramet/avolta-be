package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IncentiveSettingGetAllHandler(c *gin.Context) {
	var data []model.IncentiveSetting
	preloads := []string{}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

	paginator.Paginate(&data)
	// search, ok := c.GetQuery("search")
	// if ok {
	// 	paginator.Search = append(paginator.Search, map[string]interface{}{
	// 		"full_name": search,
	// 	})
	// }
	dataRecords, err := paginator.Paginate(&data)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponsePaginatorSuccess(c, "Data List IncentiveSetting Retrived", dataRecords.Records, dataRecords)
}

func IncentiveSettingGetOneHandler(c *gin.Context) {
	var data model.IncentiveSetting

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveSetting Retrived", data, nil)
}

func IncentiveSettingCreateHandler(c *gin.Context) {
	var data model.IncentiveSetting

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
	util.ResponseSuccess(c, "Data IncentiveSetting Created", gin.H{"last_id": data.ID}, nil)
}

func IncentiveSettingUpdateHandler(c *gin.Context) {
	var input, data model.IncentiveSetting
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
	util.ResponseSuccess(c, "Data IncentiveSetting Updated", nil, nil)
}

func IncentiveSettingDeleteHandler(c *gin.Context) {
	var input, data model.IncentiveSetting
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data IncentiveSetting Deleted", nil, nil)
}
