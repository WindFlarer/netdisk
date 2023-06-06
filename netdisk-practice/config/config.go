package config

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// 邮箱授权码
var MailPassword = "aytsiypbyufzgcia"

// 验证码长度
var CodeLength = 6

// token密码
var JwtKey = "netdisk-key"

// token过期时间
var TokenExpire = 86400

// 腾讯云
var TencentSecretKey = os.Getenv("TencentSecretKey")
var TencentSecretID = os.Getenv("TencentSecretID")
var CosBucket = "https://1-1308034550.cos.ap-nanjing.myqcloud.com"
var Dir = "netdisk-practice/"
var BucketURL = "https://1-1308034550.cos.ap-nanjing.myqcloud.com/"
