package controllers

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"netdisk-practice.com/config"
	"netdisk-practice.com/helpers"
	"netdisk-practice.com/models"
)

// 上传文件接口
func UploadFile(c *gin.Context) {
	var fileUploadRequest models.FileUploadRequest
	c.ShouldBind(&fileUploadRequest)

	// 相同文件不上传
	result := models.DB.Where("path = ? and user_name = ?", fileUploadRequest.Path, c.Request.Header.Get("UserName")).First(&models.FileBasic{})
	if result.RowsAffected > 0 { //文件已经存在
		c.JSON(400, gin.H{
			"code": 1,
			"msg":  "文件已经存在, 上传失败",
		})
		return
	}

	// 上传到腾讯cos
	uuidName := helpers.UUID()
	if suffix := path.Ext(fileUploadRequest.Path); suffix != "" {
		uuidName = uuidName + suffix
	}
	err := helpers.FileUpload(c, uuidName)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 2,
			"msg":  "上传腾讯云失败",
		})
		return
	}

	// 计算文件大小
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fileSize := file.Size

	fileBasic := models.FileBasic{
		COSPath:  config.BucketURL + config.Dir + uuidName,
		FileName: fileUploadRequest.FileName,
		UserName: c.Request.Header.Get("UserName"),
		FileSize: fileSize,
		Path:     fileUploadRequest.Path,
		IsDir:    false,
	}

	// 保存到数据库

	//创建记录
	err = models.DB.Create(&fileBasic).Error

	if err != nil {
		c.JSON(200, gin.H{
			"code": 3,
			"msg":  "上传记录失败",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "上传成功",
		})
		return
	}

}

// 下载文件接口
func DownloadFile(c *gin.Context) {
	var fileDownloadRequest models.FileDownloadRequest
	c.ShouldBind(&fileDownloadRequest)

	// 下载文件

	// 获取对应的COSpath
	cosFileName := helpers.ToCosPath(fileDownloadRequest.Path, c.Request.Header.Get("UserName"))
	if cosFileName == "" {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "文件不存在",
		})
		return
	}

	err := helpers.FileDownload(cosFileName, fileDownloadRequest.FileName, fileDownloadRequest.DownPath)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "下载成功",
	})

}

// 删除文件接口
func DeleteFile(c *gin.Context) {
	var fileDeleteFileRequest models.FileDeleteFileRequest
	c.ShouldBind(&fileDeleteFileRequest)

	// 查询文件是否存在
	cosFileName := helpers.ToCosPath(fileDeleteFileRequest.Path, c.Request.Header.Get("UserName"))
	if cosFileName == "" {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "cos文件不存在",
		})
		return
	}

	// 删除cos里面的文件
	err := helpers.FileDelete(cosFileName)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "cos删除失败",
		})
	}

	//清除数据库里面的文件, 根据path来找文件
	result := models.DB.Where("path = ? and user_name = ?", fileDeleteFileRequest.Path, c.Request.Header.Get("UserName")).Unscoped().Delete(&models.FileBasic{})
	if result.RowsAffected > 0 {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "删除成功",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "删除失败",
		})
	}

}

// 删除文件夹接口
func DeleteDir(c *gin.Context) {
	var fileDeleteDirRequest models.FileDeleteDirRequest
	c.ShouldBind(&fileDeleteDirRequest)
	var fileBasics []models.FileBasic

	// 查找原路径
	// 1. 找出所有路径包含dirName的路径, 包括文件 文件夹以及本身
	result := models.DB.Where("(path LIKE ? or (path = ?) and is_dir = 1) and user_name = ? ", fileDeleteDirRequest.Path+"/"+"%", fileDeleteDirRequest.Path, c.Request.Header.Get("UserName")).Find(&fileBasics)

	if result.RowsAffected == 0 { //没查到原来的路径
		c.JSON(200, gin.H{"code": 2, "msg": "原文件夹不存在"})
		return
	}

	// 2. 将文件 文件夹进行删除
	for i := 0; i < len(fileBasics); i++ {
		if fileBasics[i].IsDir { // 如果是文件夹, 直接删除
			if err := models.DB.Delete(&fileBasics[i]).Error; err != nil {
				if err := models.DB.Delete(&fileBasics[i]).Error; err != nil {
					c.JSON(200, gin.H{
						"code": 2,
						"msg":  "删除文件夹记录保存失败",
					})
				}
			}
		} else { // 删除文件
			//查找cosFileName, 用于删除cos中的文件
			cosFileName := helpers.ToCosPath(fileBasics[i].Path, c.Request.Header.Get("UserName"))
			if cosFileName == "" {
				c.JSON(200, gin.H{
					"code": 1,
					"msg":  "cos文件不存在",
				})
				return
			}

			// 删除cos里面的文件
			err := helpers.FileDelete(cosFileName)
			if err != nil {
				c.JSON(200, gin.H{
					"code": 1,
					"msg":  "cos删除失败",
				})
			}
			// 清除数据库记录
			if err := models.DB.Delete(&fileBasics[i]).Error; err != nil {
				c.JSON(200, gin.H{
					"code": 2,
					"msg":  "删除文件记录保存失败",
				})
			}
		}
	}

	c.JSON(200, gin.H{"code": 0, "msg": "文件夹删除成功"})
}

// 获取文件列表接口
func FileList(c *gin.Context) {
	var fileList []models.FileBasic
	if err := models.DB.Where("user_name = ?", c.Request.Header.Get("UserName")).Find(&fileList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"data": fileList,
		})
	}
}

// 移动文件接口
func MoveFile(c *gin.Context) {
	var fileMoveRequest models.FileMoveRequest
	c.ShouldBind(&fileMoveRequest)

	// 查找原路径
	var fileBasic models.FileBasic
	result := models.DB.Where("path = ? and user_name = ?", fileMoveRequest.OldPath, c.Request.Header.Get("UserName")).First(&fileBasic)
	if result.RowsAffected == 0 { //没查到原来的路径
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "原路径不存在"})
		return
	}

	// 检查新路径
	var fileBasic2 models.FileBasic
	result = models.DB.Where("path = ? and user_name = ?", fileMoveRequest.NewPath, c.Request.Header.Get("UserName")).First(&fileBasic2)
	if result.RowsAffected > 0 { // 新路径已存在
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "新路径已存在"})
		return
	}

	// 修改原路径
	fileBasic.Path = fileMoveRequest.NewPath
	err := models.DB.Save(&fileBasic).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"msg": "修改成功"})
	}
}

// 修改文件名接口
func ModifyFile(c *gin.Context) {
	var fileModifyFileRequest models.FileModifyFileRequest
	c.ShouldBind(&fileModifyFileRequest)

	// 查找原路径
	var fileBasic models.FileBasic
	result := models.DB.Where("path = ? and user_name = ? and is_dir = 0", fileModifyFileRequest.Path, c.Request.Header.Get("UserName")).First(&fileBasic)
	if result.RowsAffected == 0 { //没查到原来的路径
		c.JSON(200, gin.H{"code": 2, "msg": "原文件不存在"})
		return
	}

	newFilePath := fileModifyFileRequest.Path[:strings.LastIndex(fileModifyFileRequest.Path, "/")+1] + fileModifyFileRequest.NewFileName

	// 检查新路径
	var fileBasic2 models.FileBasic
	result = models.DB.Where("path = ? and user_name = ? and is_dir = 0", newFilePath, c.Request.Header.Get("UserName")).First(&fileBasic2)
	if result.RowsAffected > 0 { // 新路径已存在
		c.JSON(200, gin.H{"code": 1, "msg": "新路径已存在"})
		return
	}

	// 修改文件名

	fileBasic.Path = newFilePath
	fileBasic.FileName = fileModifyFileRequest.NewFileName
	err := models.DB.Save(&fileBasic).Error

	if err != nil {
		c.JSON(200, gin.H{"msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"code": 0, "msg": "修改成功"})
	}
}

// 创建文件夹接口
func CreateDir(c *gin.Context) {
	var fileCreateDir models.FileCreateDir
	c.ShouldBind(&fileCreateDir)

	// 查找原文件夹
	var fileBasic models.FileBasic
	result := models.DB.Where("path = ? and user_name = ?", fileCreateDir.Path, c.Request.Header.Get("UserName")).First(&fileBasic)
	if result.RowsAffected > 0 { //查到原来的路径
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "原文件夹已存在"})
		return
	}

	// 生成文件夹
	err := models.DB.Create(&models.FileBasic{FileName: fileCreateDir.FileName, Path: fileCreateDir.Path, IsDir: true, UserName: c.Request.Header.Get("UserName")}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"code": 0, "msg": "创建成功"})
	}
}

// 修改文件夹名
func ModifyDir(c *gin.Context) {
	var fileModifyDirRequest models.FileModifyDirRequest
	c.ShouldBind(&fileModifyDirRequest)
	var fileBasics []models.FileBasic

	// 查找原路径
	// 1. 找出所有路径包含dirName的路径
	result := models.DB.Where("(path LIKE ? or path = ?) and user_name = ? ", fileModifyDirRequest.Path+"/"+"%", fileModifyDirRequest.Path, c.Request.Header.Get("UserName")).Find(&fileBasics)

	if result.RowsAffected == 0 { //没查到原来的路径
		c.JSON(200, gin.H{"code": 2, "msg": "原文件夹不存在"})
		return
	}

	newDirPath := fileModifyDirRequest.Path[:strings.LastIndex(fileModifyDirRequest.Path, "/")+1] + fileModifyDirRequest.NewDirName

	// 检查新路径
	var fileBasic2 []models.FileBasic
	result = models.DB.Where("(path LIKE ? or path = ?)and user_name = ? and is_dir = 1", newDirPath+"/"+"%", newDirPath, c.Request.Header.Get("UserName")).Find(&fileBasic2)
	if result.RowsAffected > 0 { // 新路径已存在
		c.JSON(200, gin.H{"code": 1, "msg": "新文件夹已存在"})
		return
	}

	// 修改文件夹名(把所有路径含有oldName的都找出来)

	// 2. 修改路径
	prefix := ""
	suffix := ""
	for i := 0; i < len(fileBasics); i++ {
		prefix = fileBasics[i].Path[:len(fileModifyDirRequest.Path)-len(fileModifyDirRequest.OldDirName)]
		suffix = fileBasics[i].Path[len(fileModifyDirRequest.Path):]
		fileBasics[i].Path = prefix + fileModifyDirRequest.NewDirName + suffix
		if fileBasics[i].Path == fileModifyDirRequest.Path {
			fileBasics[i].FileName = fileModifyDirRequest.NewDirName
		}
	}

	// 3. 保存
	err := models.DB.Save(&fileBasics).Error

	if err != nil {
		c.JSON(200, gin.H{"msg": err.Error()})
	} else {
		c.JSON(200, gin.H{"code": 0, "msg": "修改成功"})
	}
}
