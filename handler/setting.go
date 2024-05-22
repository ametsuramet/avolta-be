package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SettingAutoNumberHandler(c *gin.Context) {
	var data model.Setting
	var payRoll model.PayRoll

	if err := database.DB.First(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	nextNumber := ""
	if data.PayRollAutoNumber {
		if err := database.DB.Order("created_at desc").First(&payRoll).Error; err != nil {
			nextNumber = model.GenerateInvoiceBillNumber(data, "00")
		} else {
			nextNumber = model.ExtractNumber(data, payRoll.PayRollNumber)
		}
	}

	util.ResponseSuccess(c, "Data Setting Retrived", nextNumber, nil)
}
func SettingGetOneHandler(c *gin.Context) {
	var data model.Setting
	count := int64(0)
	database.DB.Model(&data).Count(&count)

	if count == 0 {
		data = model.Setting{}
		database.DB.Create(&data)
	}

	if err := database.DB.First(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Setting Retrived", data, nil)
}

func SettingUpdateHandler(c *gin.Context) {
	var input, data model.Setting

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.First(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Updates(&input).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if !input.PayRollAutoNumber {
		if err := database.DB.Model(&data).Update("pay_roll_auto_number", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !input.IsEffectiveRateAverage {
		if err := database.DB.Model(&data).Update("is_effective_rate_average", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !input.IsGrossUp {
		if err := database.DB.Model(&data).Update("is_gross_up", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	util.ResponseSuccess(c, "Data Setting Updated", nil, nil)
}
