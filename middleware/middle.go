/**
 * @Author: LOFTER
 * @Description:
 * @File:  middle
 * @Date: 2021/2/5 3:04 下午
 */
package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"upload/util"
)

//func Cors(c *gin.Context) {
//	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
//
//	c.Next()
//}

func Cors() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie", "X-Token", "*"}
	//corsConfig.AllowHeaders = []string{"*"}
	//corsConfig.AllowOrigins = []string{
	//	"http://192.168.1.88:8080",
	//}
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	return cors.New(corsConfig)
}

//文件过滤  TODO 处理err

func FileFilterMiddle(c *gin.Context) {
	// 校验token 和 过期时间
	token := c.Request.URL.Query().Get("token")

	b, _ := util.DePwdCode(token)
	time1, _ := strconv.ParseInt(string(b), 10, 64)

	if time.Now().UnixNano() <= int64(time1) {
		c.Next()
	} else {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "禁止下载",
		})
		c.Abort()
	}

}
