/**
 * @Author: LOFTER
 * @Description:
 * @File:  chunk
 * @Date: 2021/2/5 3:00 下午
 */
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"upload/model"
	"upload/util"
)

//检测 chunk

func CheckChunk(c *gin.Context) {
	hash := c.Query("hash")
	hashPath := fmt.Sprintf("./uploadFile/%s", hash)
	chunkList := make([]string, 0)
	isExistPath, err := PathExists(hashPath)
	if err != nil {
		fmt.Println("获取hash路径错误", err)
	}

	if isExistPath {
		files, err := ioutil.ReadDir(hashPath)
		state := 0
		if err != nil {
			fmt.Println("文件读取错误", err)
		}
		for _, f := range files {

			if _, err := strconv.Atoi(f.Name()); err != nil {
				c.JSON(200, gin.H{
					"state":     600,
					"chunkList": nil,
					"msg":       "文件已经上传",
				})
				return
			}

			fileName := f.Name()
			chunkList = append(chunkList, fileName)
			fileBaseName := strings.Split(fileName, ".")[0]
			if fileBaseName == hash {
				state = 1
			}
		}
		c.JSON(200, gin.H{
			"state":     state,
			"chunkList": chunkList,
		})
	} else {
		c.JSON(200, gin.H{
			"state":     0,
			"chunkList": chunkList,
		})
	}
}

//upload

func UploadChunk(c *gin.Context) {
	fileHash := c.PostForm("hash")
	file, err := c.FormFile("file")
	hashPath := fmt.Sprintf("./uploadFile/%s", fileHash)
	if err != nil {
		fmt.Println("获取上传文件失败", err)
	}

	isExistPath, err := PathExists(hashPath)
	if err != nil {
		fmt.Println("获取hash路径错误", err)
	}

	if !isExistPath {
		_ = os.Mkdir(hashPath, os.ModePerm)
	}

	err = c.SaveUploadedFile(file, fmt.Sprintf("./uploadFile/%s/%s", fileHash, file.Filename))
	if err != nil {
		c.String(400, "0")
		fmt.Println(err)
	} else {
		chunkList := make([]string, 0)
		files, err := ioutil.ReadDir(hashPath)
		if err != nil {
			fmt.Println("文件读取错误", err)
		}
		for _, f := range files {
			chunkList = append(chunkList, f.Name())
		}

		c.JSON(200, gin.H{
			"chunkList": chunkList,
		})
	}
}

//合并

func MeagerChunk(c *gin.Context) {
	hash := c.Query("hash")
	fileName := c.Query("fileName")
	hashPath := fmt.Sprintf("./uploadFile/%s", hash)
	savePath := fmt.Sprintf("/uploadFile/%s", hash)

	isExistPath, err := PathExists(hashPath)
	if err != nil {
		fmt.Println("获取hash路径错误", err)
	}

	if !isExistPath {
		c.JSON(400, gin.H{
			"message": "文件夹不存在",
		})
		return
	}
	isExistFile, err := PathExists(hashPath + "/" + fileName)
	if err != nil {
		fmt.Println("获取hash路径文件错误", err)
	}
	fmt.Println("文件是否存在", isExistFile)
	if isExistFile {
		c.JSON(200, gin.H{
			"fileUrl": fmt.Sprintf("http://127.0.0.1.txt:9999/uploadFile/%s/%s", hash, fileName),
		})
		return
	}

	files, err := ioutil.ReadDir(hashPath)
	fmt.Println(len(files), "--------文件数---------")
	if err != nil {
		fmt.Println("合并文件读取失败", err)
	}
	completeFile, err := os.Create(hashPath + "/" + fileName)

	for k, _ := range files {
		fileBuffer, err := ioutil.ReadFile(hashPath + "/" + strconv.Itoa(k))
		if err != nil {
			fmt.Println("文件打开错误", err)
		}
		_, _ = completeFile.Write(fileBuffer)
	}
	completeFile.Close()

	//清理chunk
	fNames, err := ioutil.ReadDir(hashPath)
	if err != nil {
		fmt.Println("文件读取错误", err)
	}
	for _, f := range fNames {
		fileName := f.Name()

		if _, err := strconv.Atoi(fileName); err != nil {
			fmt.Println(err)
		} else {
			if err := os.Remove(hashPath + "/" + fileName); err != nil {
				fmt.Println(err)
			}
		}

	}

	//md5

	//计算文件大小
	fi, _ := os.Stat(hashPath + "/" + fileName)
	//写库
	fileInfo := model.FileInfo{
		FileName: fileName,
		FileSize: ByteConversionGBMBKB(fi.Size()),
		Path:     savePath + "/" + fileName,
		Hash:     hash,
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
	md5hash, _ := util.Md5sum(hashPath + "/" + fileName)
	fmt.Println(hashPath+"/"+fileName, "------name file")
	fmt.Println(md5hash, "------md5")
	fmt.Println(hash, "------hash")
	if md5hash == hash {
		c.JSON(200, gin.H{
			"url": fmt.Sprintf("http://127.0.0.1.txt:8080/uploadFile/%s/%s/?token=%s", hash, fileName, token),
			"md5": md5hash,
		})
		return
	}
	c.JSON(200, gin.H{
		"url": "",
		"msg": "上传失败 MD5不一致",
	})

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//大小转换

func ByteConversionGBMBKB(fileSize int64) string {

	if fileSize < 1024 {
		return fmt.Sprintf("%.3fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.3fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.3fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.3fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.3fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.3fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
