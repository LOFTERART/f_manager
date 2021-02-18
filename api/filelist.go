/**
 * @Author: LOFTER
 * @Description:
 * @File:  filelist
 * @Date: 2021/2/5 3:23 下午
 */
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"upload/glo"
	"upload/model"
	"upload/serializer/fileinfo"
)

//获取列表

type FileListParams struct {
	Page       int    `json:"page,omitempty"`
	Size       int    `json:"size,omitempty"`
	StartTime  string `json:"startTime,omitempty"`
	EndTime    string `json:"endTime,omitempty"`
	Repository string `json:"repository,omitempty"`
}

func GetFileList(c *gin.Context) {

	var info FileListParams

	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(200, gin.H{})
		return
	}

	var list []*model.FileInfo

	table := model.DB.Table("file_info")

	if info.Repository != "" {
		table = table.Where("file_name like ? ", "%"+info.Repository+"%")
	}

	if info.StartTime != "" && info.EndTime != "" {
		table = table.Where(" ? <= created_at AND created_at <= ?", info.StartTime, info.EndTime)
	}
	var count int64

	table.Offset((info.Page - 1) * info.Size).Limit(info.Size).Find(&list)

	model.DB.Model(&model.FileInfo{}).Where("`deleted_at` IS NULL").Count(&count)

	c.JSON(200, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"output_list": fileinfo.BuildFileSerS(list),
			"total_num":   count,
		},
	})
}

//分享

type ShareFileParams struct {
	URL string `json:"url,omitempty"`
}

func ShareFile(c *gin.Context) {

	var info ShareFileParams

	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(200, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": "http://" + glo.IPAddress + ":8080" + info.URL + "?token=smx",
	})

}

// 删除

type DelFileParams struct {
	ID       int    `json:"id,omitempty"`
	FileName string `json:"file_name"`
	FileMD5  string `json:"file_md5"` //文件md5 redis 数据删除
}

func DelFile(c *gin.Context) {

	var info DelFileParams

	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(200, gin.H{})
		return
	}
	//删除
	model.DB.Model(&model.FileInfo{}).Delete(&model.FileInfo{
		BaseInfo: model.BaseInfo{ID: info.ID},
	})
	//文件删除 TODO 暂时忽略err
	if err := os.Remove("uploadFile/" + info.FileName); err != nil {
		fmt.Println(err, "----err")
	}
	//redis 删除
	_ = glo.ROP.DelFileInfo(info.FileMD5)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})

}
