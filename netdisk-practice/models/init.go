package models

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB = InitDB()

var RDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func InitDB() *gorm.DB {
	var err error
	dsn := "root:65163326@tcp(127.0.0.1:3306)/netdisk?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func InitTable() {
	//连接数据库
	InitDB()

	//创建基本表
	DB.AutoMigrate(&UserBasic{})
	DB.AutoMigrate(&FileBasic{})
}
