package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompanyGetOneHandler(c *gin.Context) {

	getCompany, _ := c.Get("company")
	company := getCompany.(model.Company)

	if err := database.DB.First(&company).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Company Retrived", company, nil)
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
