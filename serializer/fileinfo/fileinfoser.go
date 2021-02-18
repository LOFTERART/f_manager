/**
 * @Author: LOFTER
 * @Description:
 * @File:  fileinfoser
 * @Date: 2021/2/6 12:41 下午
 */
package fileinfo

import (
	"fmt"
	"strconv"
	"time"
	"upload/model"
	"upload/util"
)

type FileInfoSer struct {
	ID        int    `json:"id" gorm:"PRIMARY_KEY"`
	CreatedAt string `json:"created_at"`
	FileName  string `json:"file_name,omitempty"`
	FileSize  string `json:"file_size,omitempty"`
	Path      string `json:"path,omitempty"`
	Hash      string `json:"hash"`
	ShareUrl  string `json:"share_url"` //分享url 默认24h
}

func ShareUrlCreateToken(hash, fileName string) string {
	//生成访问token
	expireTime := time.Now().Add(24 * time.Hour).UnixNano()
	accessTime := strconv.FormatInt(expireTime, 10)
	token, _ := util.EnPwdCode(accessTime)
	url := fmt.Sprintf("http://127.0.0.1:8080/static/uploadFile/%s?token=%s", fileName, token)
	return url
}

func FormatTime(TimeInfo time.Time) string {
	return TimeInfo.Format("2006-01-02 15:04:05")
}

func BuildFIleInfoSer(item model.FileInfo) FileInfoSer {
	return FileInfoSer{
		ID:        item.ID,
		CreatedAt: FormatTime(item.CreatedAt),
		FileName:  item.FileName,
		FileSize:  item.FileSize,
		Path:      item.Path,
		Hash:      item.Hash,
		ShareUrl:  ShareUrlCreateToken(item.Hash, item.FileName),
	}
}

// BuildFIleInfoSer

func BuildFileSerS(items []*model.FileInfo) (infos []*FileInfoSer) {
	for _, item := range items {
		fl := BuildFIleInfoSer(*item)
		infos = append(infos, &fl)
	}
	return
}
