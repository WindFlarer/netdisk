package test

import (
	"fmt"
	"net/mail"
	"testing"

	"gopkg.in/gomail.v2"
	"netdisk-practice.com/helpers"
)

func TestMail(t *testing.T) {
	// 邮件发送者的信息
	sender := mail.Address{Name: "Accepted", Address: "1309531801@qq.com"}
	// 邮件接收者的信息
	receiver := mail.Address{Name: "", Address: "2713129211@qq.com"}
	// 邮件主题
	subject := "Verification Code"
	// 邮件内容
	code := helpers.RandNum()
	body := fmt.Sprintf("Your verification code is: %s", code)

	// 邮件服务器的配置
	host := "smtp.qq.com"           // QQ邮箱SMTP服务器的地址
	port := 465                     // QQ邮箱SMTP服务器的端口号
	username := "1309531801@qq.com" // QQ邮箱的登录用户名
	password := "0919SHENGRI,"      // QQ邮箱的登录密码

	// 创建邮件对象
	m := gomail.NewMessage()
	m.SetHeader("From", sender.String())
	m.SetHeader("To", receiver.String())
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// 创建邮件发送器
	d := gomail.NewDialer(host, port, username, password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}
