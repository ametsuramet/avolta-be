package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PayRollItemGetAllHandler(c *gin.Context) {
	var data []model.PayRollItem
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

	util.ResponsePaginatorSuccess(c, "Data List PayRollItem Retrived", dataRecords.Records, dataRecords)
}

func PayRollItemGetOneHandler(c *gin.Context) {
	var data model.PayRollItem

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRollItem Retrived", data, nil)
}

func PayRollItemCreateHandler(c *gin.Context) {
	var input model.PayRollItemReq
	var data model.PayRollItem
	var payRoll model.PayRoll

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	// getUser, _ := c.Get("user")
	// user := getUser.(model.User)

	// data.AuthorID = user.ID
	getCompany, _ := c.Get("company")
	company := getCompany.(model.Company)
	data.CompanyID = company.ID

	if err := database.DB.Find(&payRoll, "id = ?", input.PayRollID).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Create(&model.PayRollItem{
		ItemType:        input.ItemType,
		Title:           input.Title,
		Notes:           input.Notes,
		IsDefault:       input.IsDefault,
		IsDeductible:    input.IsDeductible,
		IsTax:           input.IsTax,
		TaxAutoCount:    input.TaxAutoCount,
		IsTaxCost:       input.IsTaxCost,
		IsTaxAllowance:  input.IsTaxAllowance,
		Amount:          input.Amount,
		PayRollID:       input.PayRollID,
		ReimbursementID: input.ReimbursementID,
	}).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	payRoll.CountTax()
	util.ResponseSuccess(c, "Data PayRollItem Created", gin.H{"last_id": data.ID}, nil)
}

func PayRollItemUpdateHandler(c *gin.Context) {
	var input model.PayRollItemReq
	var data model.PayRollItem
	var payRoll model.PayRoll
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Find(&payRoll, "id = ?", input.PayRollID).Error; err != nil {
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
	if input.Amount == 0 {
		if err := database.DB.Model(&data).Update("amount", 0).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	payRoll.CountTax()
	util.ResponseSuccess(c, "Data PayRollItem Updated", nil, nil)
}

func PayRollItemDeleteHandler(c *gin.Context) {
	var data model.PayRollItem
	id := c.Params.ByName("id")

	if err := database.DB.Model(&data).Unscoped().Delete(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRollItem Deleted", nil, nil)
}
