/**
 * @Author: LOFTER
 * @Description:
 * @File:  model
 * @Date: 2021/2/5 1.txt:48 下午
 */

package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

//Initialized  连接mysql
func Initialized() {
	// 从本地读取环境变量
	DB, _ = gorm.Open(sqlite.Open("fileManager.db"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})

	_ = DB.AutoMigrate(&FileInfo{}, &Admin{})
}
