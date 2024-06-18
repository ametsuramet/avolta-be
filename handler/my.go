package handler

import (
	"avolta/config"
	"avolta/model"
	"avolta/util"
	"fmt"

	"github.com/gin-gonic/gin"
)

func MyController(c *gin.Context) {
	getUser, _ := c.Get("user")
	user := getUser.(model.User)
	getEmployee, ok := c.Get("employee")
	getRoles, _ := c.Get("roles")
	roles := getRoles.([]model.Role)
	if ok {
		user.Employee = getEmployee.(model.Employee)
		if user.Employee.Picture.Valid {
			user.Avatar = fmt.Sprintf("%s/%s", config.App.Server.BaseURL, user.Employee.Picture.String)
		}
	}
	if len(roles) > 0 {
		for _, v := range roles {
			if v.IsSuperAdmin {
				v.GetPermissions()
				for _, v := range v.Permissions {
					user.Permissions = append(user.Permissions, v.Key)
				}
			}
		}
	}
	util.ResponseSuccess(c, "Profile Retrieved", user, nil)
}
