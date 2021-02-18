package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"upload/api"
	"upload/glo"
	"upload/middleware"
	"upload/model"
	"upload/tool"
	"upload/util"
)

func init() {
	model.Initialized()
}

func init() {
	tool.InitLogger()
	glo.ROP.RedisInit()
	glo.IPAddress, _ = util.GetIPV4()
}

func main() {
	//admin create
	CreateAdmin()
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.Use(middleware.Cors())
	admin := router.Group("/adminService")
	{
		//获取file
		admin.POST("/get_file_list", api.GetFileList)
		//分享
		admin.POST("/share", api.ShareFile)
		//删除一个
		admin.POST("/del_one", api.DelFile)
	}

	user := router.Group("user")
	{

		//注册admin用户
		user.POST("/v1/createAdmin", api.CreateAdminUser)
		//登录admin用户
		user.POST("/v1/loginAdmin", api.LoginAdminUser)
		//获取后台登录用户信息
		user.POST("/v1/getAdminInfo", api.GetAdminUserInfo)
		//登出
		user.POST("/v1/logout", api.LoginOut)
	}

	//Chunk := router.Group("/upload")
	//{
	//
	//	Chunk.GET("/checkChunk", api.CheckChunk)
	//
	//	Chunk.POST("/uploadChunk", api.UploadChunk)
	//
	//	Chunk.GET("/meagerChunk", api.MeagerChunk)
	//}
	//使用redis 记录 chunk 信息
	Chunk := router.Group("/upload/redis")
	{

		Chunk.GET("/checkChunk", api.CheckChunkNew)

		Chunk.POST("/uploadChunk", api.UploadChunkNew)

	}

	Static := router.Group("/static")
	{
		Static.Use(middleware.FileFilterMiddle)
		Static.StaticFS("/uploadFile", http.Dir("./uploadFile"))
	}

	_ = router.Run(":8080")
}

//生成管理员

func CreateAdmin() {
	user := model.Admin{
		Name:     "admin",
		Password: "111111",
		Tokens:   "smx",
		Roles:    "admin",
	}
	model.DB.Create(&user)
}
