/**
 * @Author: LOFTER
 * @Description:
 * @File:  newchunk
 * @Date: 2021/2/8 12:50 上午
 */
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
	"upload/chunk"
	"upload/glo"
	"upload/model"
	"upload/util"
)

type JsonRes struct {
	Code    int         `json:"code"`
	Success bool        `json:"success,omitempty"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

//检测 chunk

func CheckChunkNew(c *gin.Context) {

	form := c.Request.URL.Query()
	fileMd5 := form["fileMd5"][0]
	// 查文件是否存在
	f := glo.ROP.GetFileinfo(fileMd5)
	if len(f) > 0 { // 文件不存在直接返回
		c.JSON(200, &JsonRes{0, true, "文件存在", f})
		return
	}

	// 文件不存在就检测是否有在上传当中的 断点续传
	mem := glo.ROP.GetMemWithScore(fileMd5)
	if len(mem) > 0 {
		c.JSON(200, &JsonRes{301, true, "断点上传", mem})
		return
	}

	c.JSON(200, &JsonRes{401, true, "文件不存在", nil})

}

//upload

func UploadChunkNew(c *gin.Context) {
	//设置内存大小
	_ = c.Request.ParseMultipartForm(32 << 20)
	//获取上传的第一个文件
	filename := c.Request.Form.Get("filename")                                      //文件名
	filetype := c.Request.Form.Get("type")                                          //文件类型
	file, _, err := c.Request.FormFile("file")                                      //file
	fileMd5 := c.Request.Form.Get("fileMd5")                                        //file md5
	fileSize := c.Request.Form.Get("fileSize")                                      //file size
	lastModified, _ := strconv.ParseInt(c.Request.Form.Get("lastModified"), 10, 64) //lastModified
	idx, _ := strconv.Atoi(c.Request.Form.Get("index"))                             //current chunk index
	pieceMd5 := c.Request.Form.Get("chunkMd5")                                      // 分片md5
	chunks, _ := strconv.Atoi(c.Request.Form.Get("chunks"))                         //total chunks
	defer file.Close()
	if err != nil {
		c.JSON(200, &JsonRes{0, true, "文件上传出错", err.Error()})
		return
	}
	// 文件如果存在分片也是没有必要的
	info := glo.ROP.GetFileinfo(fileMd5)
	if len(info) > 0 {
		c.JSON(200, &JsonRes{0, true, "文件存在", info})
		return
	}

	chunk.PieceSave(file, pieceMd5, fileMd5) // 保存分片
	glo.ROP.ChunkAdd(idx, fileMd5, pieceMd5)
	// 上传完并且没有正在合并的则可以执行合并
	if chunks == glo.ROP.ChunkIsFull(fileMd5) && glo.ROP.ISMerging(fileMd5) == 0 {
		glo.ROP.Merging(fileMd5) // 合并先加锁
		all := glo.ROP.GetMem(fileMd5)
		_, md5Name, filepath := glo.FD.MergeFile(fileMd5, filename, all, chunks, lastModified)
		glo.ROP.FileInfo(fileMd5, filename, filepath, filetype) // 保存文件信息
		glo.ROP.DelMerging(fileMd5)                             // 文件信息存储到map中后就可以解除锁了

		// 文件md5校验
		if md5Name != fileMd5 {
			c.JSON(200, &JsonRes{400, false, "文件不一致", nil})
			return
		}
		glo.ROP.ClearSet(fileMd5) // 满了则清除集合 合并文件

		//删除chunk
		_ = os.RemoveAll("uploadFile/chunks/" + fileMd5)

		//写库
		size, _ := strconv.ParseInt(fileSize, 10, 64)
		fileInfo := model.FileInfo{
			FileName: filename,
			FileSize: ByteConversionGBMBKB(size),
			Path:     "uploadFile" + "/" + filename,
			Hash:     fileMd5,
		}
		model.DB.Create(&fileInfo)

		//生成访问token
		expireTime := time.Now().Add(24 * time.Hour).UnixNano()

		accessTime := strconv.FormatInt(expireTime, 10)

		token, err := util.EnPwdCode(accessTime)
		if err != nil {
			c.JSON(200, "fail")
			return
		}
		url := fmt.Sprintf("http://127.0.0.1.txt:8080/uploadFile/%s?token=%s", filename, token)

		c.JSON(200, &JsonRes{0, true, "文件合并完成", map[string]string{"url": url}})
		return
	}
	c.JSON(200, &JsonRes{0, true, "分片上传完成", idx})

}
