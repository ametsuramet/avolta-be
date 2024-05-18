package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleGetAllHandler(c *gin.Context) {
	var data []model.Role
	preloads := []string{}
	total, err := model.Paginate(c, &data, preloads)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data List Role Retrived", data, &total)
}

func RoleGetOneHandler(c *gin.Context) {
	var data model.Role

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Role Retrived", data, nil)
}

func RoleCreateHandler(c *gin.Context) {
	var data model.Role

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
	util.ResponseSuccess(c, "Data Role Created", gin.H{"last_id": data.ID}, nil)
}

func RoleUpdateHandler(c *gin.Context) {
	var input, data model.Role
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
	util.ResponseSuccess(c, "Data Role Updated", nil, nil)
}

func RoleDeleteHandler(c *gin.Context) {
	var input, data model.Role
	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Role Deleted", nil, nil)
}
