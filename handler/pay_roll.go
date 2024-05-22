package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PayRollGetAllHandler(c *gin.Context) {
	var data []model.PayRoll
	preloads := []string{"Employee"}
	total, err := model.Paginate(c, &data, preloads)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data List PayRoll Retrived", data, &total)
}

func PayRollGetOneHandler(c *gin.Context) {
	var data model.PayRoll

	id := c.Params.ByName("id")

	if err := database.DB.Preload("Employee").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	// data.CountTax()
	util.ResponseSuccess(c, "Data PayRoll Retrived", data, nil)
}

func PayRollCreateHandler(c *gin.Context) {
	var data model.PayRoll

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Create(&data).Error; err != nil {
			return err
		}
		data.GetEmployee()
		data.CreateDefaultItems(c)
		data.CountTax()
		return nil
	}); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRoll Created", gin.H{"last_id": data.ID}, nil)
}

func PayRollUpdateHandler(c *gin.Context) {
	var input, data model.PayRoll
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
			return err
		}
		database.DB.Model(&data).Update("is_gross_up", input.IsGrossUp)
		database.DB.Model(&data).Update("is_effective_rate_average", input.IsEffectiveRateAverage)

		data.IsGrossUp = input.IsGrossUp
		data.IsEffectiveRateAverage = input.IsEffectiveRateAverage

		data.CountTax()
		return nil
	}); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data PayRoll Updated", nil, nil)
}

func PayRollDeleteHandler(c *gin.Context) {
	var input, data model.PayRoll
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRoll Deleted", nil, nil)
}
