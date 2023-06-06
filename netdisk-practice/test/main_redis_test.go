package test

import (
	"crypto/tls"
	"net/smtp"
	"testing"
	"time"

	"github.com/jordan-wright/email"
	"netdisk-practice.com/config"
	"netdisk-practice.com/helpers"
)

func TestMailRedis(t *testing.T) {
	e := email.NewEmail()
	e.From = "Accepted <1309531801@qq.com>" //发件人
	e.To = []string{"2713129211@qq.com"}    //收件人
	e.Subject = "测试"                        //标题
	code := helpers.RandNum()
	e.HTML = []byte("验证码为:" + code) //测试代码

	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "1309531801@qq.com", config.MailPassword, "smtp.qq.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		t.Fatal(err)
	}

	err = rdb.Set(ctx, "email", code, 300*time.Second).Err()
	if err != nil {
		t.Error(err)
	}
}
