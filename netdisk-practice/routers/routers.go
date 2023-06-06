package routers

import "github.com/gin-gonic/gin"

func InitRouters(r *gin.Engine) {
	//初始化用户路由
	InitUserRouters(r)
	InitFileRouters(r)
}
