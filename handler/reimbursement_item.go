package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReimbursementItemGetAllHandler(c *gin.Context) {
	var data []model.ReimbursementItem
	preloads := []string{}
	paginator := util.NewPaginator(c)
	paginator.Preloads = preloads

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

	util.ResponsePaginatorSuccess(c, "Data List ReimbursementItem Retrived", dataRecords.Records, dataRecords)
}

func ReimbursementItemGetOneHandler(c *gin.Context) {
	var data model.ReimbursementItem

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data ReimbursementItem Retrived", data, nil)
}

func ReimbursementItemCreateHandler(c *gin.Context) {
	var input model.ReimbursementItemReq
	var reimbursement model.Reimbursement

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// getUser, _ := c.Get("user")
	// user := getUser.(model.User)

	var data = model.ReimbursementItem{
		Amount:          input.Amount,
		Notes:           input.Notes,
		ReimbursementID: input.ReimbursementID,
		Files:           input.Files,
	}

	if err := database.DB.Create(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Preload("Items").Find(&reimbursement, "id = ?", input.ReimbursementID).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data ReimbursementItem Created", gin.H{"last_id": data.ID}, nil)
}

func ReimbursementItemUpdateHandler(c *gin.Context) {
	var input, data model.ReimbursementItem
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
	util.ResponseSuccess(c, "Data ReimbursementItem Updated", nil, nil)
}

func ReimbursementItemDeleteHandler(c *gin.Context) {
	var input, data model.ReimbursementItem
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data ReimbursementItem Deleted", nil, nil)
}
