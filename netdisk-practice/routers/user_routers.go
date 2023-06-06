package routers

import (
	"github.com/gin-gonic/gin"
	"netdisk-practice.com/controllers"
)

func InitUserRouters(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/login", controllers.LoginController)       // 用户登录接口
		user.POST("/register", controllers.RegisterController) // 用户注册接口
		user.POST("/code", controllers.CodeController)         // 用户接收验证码接口
	}
}
