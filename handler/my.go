package handler

import (
	"avolta/model"
	"avolta/util"

	"github.com/gin-gonic/gin"
)

func MyController(c *gin.Context) {
	getUser, _ := c.Get("user")
	user := getUser.(model.User)
	util.ResponseSuccess(c, "Profile Retrieved", user, nil)
}
