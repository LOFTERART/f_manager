/**
 * @Author: LOFTER
 * @Description:
 * @File:  admin_service
 * @Date: 2021/2/6 6:23 下午
 */
package admin

import (
	"github.com/gin-gonic/gin"
	"upload/model"
	"upload/serializer"
	"upload/util"
)

type Admin struct {
	Name     string `form:"username" json:"username" `
	Password string `form:"password" json:"password" `
	Tokens   string `form:"tokens" json:"tokens" `
}

func (item *Admin) CreateAdmin() serializer.Response {

	admin := model.Admin{
		Name:     item.Name,
		Password: item.Password,
		Tokens:   util.RandStringRunes(10),
	}

	//查询数据库是否存在

	model.DB.Where("name = ?", item.Name).First(&admin)

	if admin.ID > 0 {
		return serializer.Response{
			Code:  20001,
			Data:  1,
			Msg:   "用户名已存在 重新输入用户名",
			Error: "",
		}
	} else {
		model.DB.Where(model.Admin{Name: item.Name}).FirstOrCreate(&admin)

		return serializer.Response{
			Code: 0,
			Data: serializer.BuildAdminSerializer(admin),
		}

	}

}

//登录
func (item *Admin) LoginAdmin() serializer.Response {

	admin := model.Admin{
		Name:     item.Name,
		Password: item.Password,
	}

	//查询数据库是否存在

	model.DB.Model(&model.Admin{}).Where("name = ? AND password =? ", item.Name, item.Password).First(&admin)

	if admin.ID == 0 {
		return serializer.Response{
			Code: 20001,
			Msg:  "请检查用户名或者密码",
		}
	} else {
		return serializer.Response{
			Code: 0,
			Data: gin.H{
				"token": admin.Tokens,
			},
		}

	}

}

//获取用户信息
func (item *Admin) GetAdmin() serializer.Response {
	var admin model.Admin
	//查询数据库是否存在
	model.DB.Where("tokens = ? ", item.Tokens).First(&admin)
	//设置全局路径变量
	return serializer.Response{
		Code: 0,
		Data: serializer.BuildAdminSerializer(admin),
	}

}
