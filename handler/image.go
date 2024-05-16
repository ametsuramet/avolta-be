package handler

import (
	"avolta/database"
	"avolta/model"
	"avolta/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ImageGetAllHandler(c *gin.Context) {
	var data []model.Image
	preloads := []string{}
	total, err := model.Paginate(c, &data, preloads)
	if err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data List Image Retrived", data, &total)
}

func ImageGetOneHandler(c *gin.Context) {
	var data model.Image

	id := c.Params.ByName("id")

	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Image Retrived", data, nil)
}

func ImageCreateHandler(c *gin.Context) {
	var data model.Image

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
	util.ResponseSuccess(c, "Data Image Created", gin.H{"last_id": data.ID}, nil)
}

func ImageUpdateHandler(c *gin.Context) {
	var input, data model.Image
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
	util.ResponseSuccess(c, "Data Image Updated", nil, nil)
}

func ImageDeleteHandler(c *gin.Context) {
	var input, data model.Image
	id := c.Params.ByName("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Find(&data, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := database.DB.Model(&data).Delete(&input, "id = ?", id).Error; err != nil {
		util.ResponseFail(c, http.StatusBadRequest, err.Error())
		return
	}
	util.ResponseSuccess(c, "Data Image Deleted", nil, nil)
}
