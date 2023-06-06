package test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
	"netdisk-practice.com/config"
)

func TestDownload(t *testing.T) {
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

	key := "netdisk-practice/pic1.jpg" // 前面不需要带斜杠/

	// opt 可选，无特殊设置可设为 nil
	// 1. 从响应体中获取对象
	resp, err := client.Object.Get(context.Background(), key, nil)
	if err != nil {
		panic(err)
	}
	ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// 2. 下载对象到本地文件
	_, err = client.Object.GetToFile(context.Background(), key, "test.jpg", nil)
	if err != nil {
		panic(err)
	}

}
