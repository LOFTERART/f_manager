/**
 * @Author: LOFTER
 * @Description:
 * @File:  baseinfo
 * @Date: 2021/2/5 1.txt:49 下午
 */
package model

import (
	"gorm.io/gorm"
	"time"
)

//基础字段common
type BaseInfo struct {
	ID        int            `json:"id" gorm:"PRIMARY_KEY"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
