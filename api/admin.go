/**
 * @Author: LOFTER
 * @Description:
 * @File:  admin
 * @Date: 2021/2/6 6:20 下午
 */
package api

import (
	"github.com/gin-gonic/gin"
	"upload/serializer"
	"upload/service/admin"
)

func CreateAdminUser(c *gin.Context) {
	var userAdmin admin.Admin

	if err := c.ShouldBind(&userAdmin); err != nil {
		c.JSON(201, err)
	} else {
		res := userAdmin.CreateAdmin()
		c.JSON(200, &res)
	}

}

func LoginAdminUser(c *gin.Context) {
	var userAdmin admin.Admin

	if err := c.ShouldBind(&userAdmin); err != nil {
		c.JSON(201, err)
	} else {
		res := userAdmin.LoginAdmin()
		c.JSON(200, &res)
	}

}

func LoginOut(c *gin.Context) {
	c.JSON(200, &serializer.Response{Code: 0, Msg: "ok"})

}

func GetAdminUserInfo(c *gin.Context) {
	var userAdmin admin.Admin
	if err := c.ShouldBind(&userAdmin); err != nil {
		c.JSON(201, err)
	} else {
		res := userAdmin.GetAdmin()
		c.JSON(200, &res)
	}

}
