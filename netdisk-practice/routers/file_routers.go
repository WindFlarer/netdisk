package routers

import (
	"github.com/gin-gonic/gin"
	"netdisk-practice.com/controllers"
	"netdisk-practice.com/middlewares"
)

func InitFileRouters(r *gin.Engine) {
	// 文件接口都需要登录才能使用
	file := r.Group("/file")
	file.Use(middlewares.AuthMiddleware())
	{
		file.POST("/upload", controllers.UploadFile)       // 文件上传
		file.POST("/download", controllers.DownloadFile)   // 文件下载
		file.DELETE("/deleteFile", controllers.DeleteFile) // 文件删除
		file.DELETE("/deleteDir", controllers.DeleteDir)   // 文件夹删除
		file.GET("/fileList", controllers.FileList)        // 文件列表
		file.PUT("/move", controllers.MoveFile)            // 移动文件
		file.PUT("/modifyFile", controllers.ModifyFile)    // 修改文件名
		file.POST("/createDir", controllers.CreateDir)     // 创建文件夹
		file.PUT("/modifyDir", controllers.ModifyDir)      // 修改文件夹名
	}
}
