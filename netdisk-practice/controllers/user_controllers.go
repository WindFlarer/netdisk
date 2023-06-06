package controllers

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"netdisk-practice.com/config"
	"netdisk-practice.com/helpers"
	"netdisk-practice.com/models"
)

// 登录
func LoginController(c *gin.Context) {
	//初始化
	var userBasic models.UserBasic
	c.ShouldBind(&userBasic)

	//查询是否存在当前用户
	nameResult := models.DB.Where("user_name = ? ", userBasic.UserName).First(&models.UserBasic{})

	//判断是否存在
	if nameResult.RowsAffected > 0 { //如果存在
		if passwordResult := models.DB.Where("user_name =? AND password =?", userBasic.UserName, userBasic.Password).First(&models.UserBasic{}); passwordResult.RowsAffected > 0 { //密码正确
			//生成token
			token, _ := helpers.GenerateToken(userBasic.ID, userBasic.UserName, config.TokenExpire)
			c.JSON(200, gin.H{
				"code":  0,
				"msg":   "登陆成功",
				"token": token,
			})
		} else { //密码错误
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "密码错误",
			})
		}
	} else { //不存在
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "该用户不存在",
		})
	}
}

// 注册
func RegisterController(c *gin.Context) {
	//初始化
	var userRegisterRequest models.UserRegisterRequest
	c.ShouldBind(&userRegisterRequest)

	//查询是否存在当前用户
	nameResult := models.DB.Where("user_name = ? ", userRegisterRequest.UserName).First(&models.UserBasic{})

	//判断是否存在
	if nameResult.RowsAffected > 0 { //如果存在
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "用户名已存在",
		})
	} else { //不存在

		// 验证码查询
		val, _ := models.RDB.Get(context.Background(), userRegisterRequest.Email).Result()

		if val == userRegisterRequest.Code { //验证码正确
			//创建用户
			userBasic := models.UserBasic{
				UserName: userRegisterRequest.UserName,
				Password: userRegisterRequest.Password,
				Email:    userRegisterRequest.Email,
				Phone:    userRegisterRequest.Phone,
			}

			// 保存用户
			models.DB.Create(&userBasic)
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "注册成功",
			})
		} else {
			c.JSON(200, gin.H{
				"code": 2,
				"msg":  "邮箱验证码错误",
			})
		}
	}
}

// 用户验证码
func CodeController(c *gin.Context) {
	// 初始化
	var userMailRequest models.UserMailRequest
	c.ShouldBind(&userMailRequest)

	//查询email是否被注册过
	var userBasic models.UserBasic
	result := models.DB.Where("email = ?", userMailRequest.Email).First(&userBasic)
	if result.RowsAffected > 0 {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "邮箱已被注册",
		})
		return
	}

	// 调用helper中的函数发送验证码
	code, err := helpers.SendMail(userMailRequest.Email)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "验证码发送失败",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "验证码发送成功",
		})
	}
	// 保存本地
	err = models.RDB.Set(context.Background(), userMailRequest.Email, code, 300*time.Second).Err()

	if err != nil {
		log.Fatal(err)
	}

}
