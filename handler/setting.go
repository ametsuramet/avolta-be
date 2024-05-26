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
			nextNumber = model.GenerateInvoiceBillNumber(model.AutoNumber{AutoNumber: data.PayRollAutoNumber, AutoFormat: data.PayRollAutoFormat, StaticCharacter: data.PayRollStaticCharacter, AutoNumberCharacterLength: data.PayRollAutoNumberCharacterLength}, "00")
		} else {
			nextNumber = model.ExtractNumber(model.AutoNumber{AutoNumber: data.PayRollAutoNumber, AutoFormat: data.PayRollAutoFormat, StaticCharacter: data.PayRollStaticCharacter, AutoNumberCharacterLength: data.PayRollAutoNumberCharacterLength}, payRoll.PayRollNumber)
		}
	}

	util.ResponseSuccess(c, "Data Setting Retrived", nextNumber, nil)
}
func SettingIncentiveAutoNumberHandler(c *gin.Context) {
	var data model.Setting
	var incentive model.IncentiveReport

	if err := database.DB.First(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	nextNumber := ""
	if data.IncentiveAutoNumber {
		if err := database.DB.Order("created_at desc").First(&incentive).Error; err != nil {
			nextNumber = model.GenerateInvoiceBillNumber(model.AutoNumber{AutoNumber: data.IncentiveAutoNumber, AutoFormat: data.IncentiveAutoFormat, StaticCharacter: data.IncentiveStaticCharacter, AutoNumberCharacterLength: data.IncentiveAutoNumberCharacterLength}, "00")
		} else {
			nextNumber = model.ExtractNumber(model.AutoNumber{AutoNumber: data.IncentiveAutoNumber, AutoFormat: data.IncentiveAutoFormat, StaticCharacter: data.IncentiveStaticCharacter, AutoNumberCharacterLength: data.IncentiveAutoNumberCharacterLength}, incentive.ReportNumber)
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
	if !input.BpjsKes {
		if err := database.DB.Model(&data).Update("bpjs_kes", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !input.BpjsTkJht {
		if err := database.DB.Model(&data).Update("bpjs_tk_jht", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !input.BpjsTkJkm {
		if err := database.DB.Model(&data).Update("bpjs_tk_jkm", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !input.BpjsTkJp {
		if err := database.DB.Model(&data).Update("bpjs_tk_jp", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	if !input.BpjsTkJkk {
		if err := database.DB.Model(&data).Update("bpjs_tk_jkk", false).Error; err != nil {
			util.ResponseFail(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	util.ResponseSuccess(c, "Data Setting Updated", nil, nil)
}
