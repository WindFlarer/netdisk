package helpers

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"

	uuid "github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"netdisk-practice.com/config"
	"netdisk-practice.com/models"
)

// 随机生成验证码
func RandNum() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < config.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

// 发送邮箱验证码
func SendMail(to string) (string, error) {
	e := email.NewEmail()
	e.From = "Accepted <1309531801@qq.com>" //发件人
	e.To = []string{to}                     //收件人
	e.Subject = "测试"                        //标题
	code := RandNum()
	e.HTML = []byte("验证码为:" + code) //测试代码

	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "1309531801@qq.com", config.MailPassword, "smtp.qq.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		log.Fatal(err)
	}
	return code, err
}

// func SendMail(to string) (string, error) {
// 	// 邮件发送者的信息
// 	sender := mail.Address{Name: "Accepted", Address: "1309531801@qq.com"}
// 	// 邮件接收者的信息
// 	receiver := mail.Address{Name: "", Address: to}
// 	// 邮件主题
// 	subject := "Verification Code"
// 	// 邮件内容
// 	code := RandNum()
// 	body := fmt.Sprintf("Your verification code is: %s", code)

// 	// 邮件服务器的配置
// 	host := "smtp.qq.com"           // QQ邮箱SMTP服务器的地址
// 	port := 465                     // QQ邮箱SMTP服务器的端口号
// 	username := "1309531801@qq.com" // QQ邮箱的登录用户名
// 	password := "0919SHENGRI,"      // QQ邮箱的登录密码

// 	// 创建邮件对象
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", sender.String())
// 	m.SetHeader("To", receiver.String())
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/plain", body)

// 	// 创建邮件发送器
// 	d := gomail.NewDialer(host, port, username, password)

// 	// 发送邮件
// 	if err := d.DialAndSend(m); err != nil {
// 		return "", err
// 	}

// 	return code, nil
// }

// 生成令牌
func GenerateToken(id uint, userName string, second int) (string, error) {
	//生成载体
	uc := config.UserClaim{
		ID:       id,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(second) * time.Second).Unix(),
		},
	}

	//生成令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(config.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析令牌
func AnalyzeToken(token string) (*config.UserClaim, error) {
	uc := new(config.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, err
}

// 获取uuid
func UUID() string {
	return uuid.NewV4().String()
}

// 转化成COSPath

func ToCosPath(path, userName string) string {
	var fileBasic models.FileBasic
	result := models.DB.Where("path = ? and user_name = ?", path, userName).First(&fileBasic)
	if result.RowsAffected == 0 {
		return ""
	}

	cosFileName := fileBasic.COSPath[strings.LastIndex(fileBasic.COSPath, "/")+1:]

	return cosFileName

}

// 文件上传
func FileUpload(c *gin.Context, uuidName string) error {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(config.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: config.TencentSecretID,
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: config.TencentSecretKey,
		},
	})

	// 文件在cos中的路径
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return err
	}
	key := config.Dir + uuidName

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// 文件下载
func FileDownload(cosFileName, fileName, downPath string) error {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(config.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: config.TencentSecretID,
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: config.TencentSecretKey,
		},
	})

	key := config.Dir + cosFileName // 前面不需要带斜杠/

	// opt 可选，无特殊设置可设为 nil
	// 1. 从响应体中获取对象
	resp, err := client.Object.Get(context.Background(), key, nil)
	if err != nil {
		return err
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 2. 下载对象到本地文件

	// 检查文件夹是否存在, 如果不存在就创建文件夹
	// 检查目录是否存在
	if _, err := os.Stat(downPath); os.IsNotExist(err) {
		// 创建目录
		err := os.MkdirAll(downPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	_, err = client.Object.GetToFile(context.Background(), key, downPath+"/"+fileName, nil)
	if err != nil {
		return err
	}

	return nil
}

// 文件删除
func FileDelete(cosFileName string) error {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(config.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: config.TencentSecretID,
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: config.TencentSecretKey,
		},
	})

	key := config.Dir + cosFileName // 前面不需要带斜杠/

	_, err := client.Object.Delete(context.Background(), key)
	if err != nil {
		return err
	}

	return nil
}
