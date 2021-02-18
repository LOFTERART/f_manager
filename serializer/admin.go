/**
 * @Author: LOFTER
 * @Description:
 * @File:  admin
 * @Date: 2021/2/6 6:21 下午
 */
package serializer

import (
	"upload/model"
)

type Admin struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Tokens string `json:"token"`
	Roles  string `json:"roles"`
	Avatar string `json:"avatar"`
}

func BuildAdminSerializer(item model.Admin) *Admin {
	return &Admin{
		Id:     item.ID,
		Name:   item.Name,
		Tokens: item.Tokens,
		Roles:  item.Roles,
		Avatar: item.Avatar,
	}
}

func BuildAdminSerializers(item []*model.Admin) (items []*Admin) {

	for _, v := range item {
		items = append(items, BuildAdminSerializer(*v))
	}
	return

}
