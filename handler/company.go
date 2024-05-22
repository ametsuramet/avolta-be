package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompanyGetOneHandler(c *gin.Context) {
	var data model.Company
	count := int64(0)
	database.DB.Model(&data).Count(&count)

	if count == 0 {
		data = model.Company{
			Name:    "Nama Perusahaan",
			Address: "Alamat Perusahaan",
		}
		database.DB.Create(&data)
	}

	if err := database.DB.First(&data).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Company Retrived", data, nil)
}

func CompanyUpdateHandler(c *gin.Context) {
	var input, data model.Company

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
	util.ResponseSuccess(c, "Data Company Updated", nil, nil)
}
