package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"netdisk-practice.com/models"
	"netdisk-practice.com/routers"
)

func main() {
	//数据库初始化
	models.InitTable()

	//路由初始化
	r := gin.Default()

	// 跨域访问
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	//加载路由
	routers.InitRouters(r)

	//启动服务
	r.Run()

}
