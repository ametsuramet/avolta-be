package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/service"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PayRollGetAllHandler(c *gin.Context) {
	var data []model.PayRoll
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

	// for i := range data {
	// 	data[i].GetTransactions()
	// 	data[i].GetPayableTransactions()
	// }

	util.ResponsePaginatorSuccess(c, "Data List PayRoll Retrived", dataRecords.Records, dataRecords)
}

func PayRollGetOneHandler(c *gin.Context) {
	var data model.PayRoll

	id := c.Params.ByName("id")

	if err := database.DB.Preload("Employee").Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	data.GetTransactions()
	data.GetPayableTransactions()
	// data.CountTax()
	util.ResponseSuccess(c, "Data PayRoll Retrived", data, nil)
}

func PayRollCreateHandler(c *gin.Context) {
	var data model.PayRoll
	var setting model.Setting

	if err := c.ShouldBindJSON(&data); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.First(&setting).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	data.BpjsSetting = service.InitBPJS()
	data.BpjsSetting.BpjsKesEnabled = setting.BpjsKes
	data.BpjsSetting.BpjsTkJhtEnabled = setting.BpjsTkJht
	data.BpjsSetting.BpjsTkJkmEnabled = setting.BpjsTkJkm
	data.BpjsSetting.BpjsTkJpEnabled = setting.BpjsTkJp
	data.BpjsSetting.BpjsTkJkkEnabled = setting.BpjsTkJkk

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Create(&data).Error; err != nil {
			return err
		}
		data.GetEmployee()
		if err := data.CreateDefaultItems(c); err != nil {
			return err
		}
		if err := data.CountTax(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRoll Created", gin.H{"last_id": data.ID}, nil)
}

func PayRollProcessHandler(c *gin.Context) {
	var data model.PayRoll
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := data.RunPayRoll(c); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}

	util.ResponseSuccess(c, "Data PayRoll Processed", nil, nil)

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

func PayRollPaymentHandler(c *gin.Context) {
	var data model.PayRoll
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := data.Payment(c); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data PayRoll Paid", nil, nil)
}
