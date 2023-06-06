package test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
	"netdisk-practice.com/config"
)

func TestMove(t *testing.T) {
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

	source := "netdisk-practice/pic.jpg"
	f := strings.NewReader("file")
	// 上传文件
	_, err := client.Object.Put(context.Background(), source, f, nil)
	if err != nil {
		panic(err)
	}
	// 移动对象
	dest := "cloud-disk/pic.jpg"
	sourceURL := fmt.Sprintf("%s/%s", u.Host, source)
	_, _, err = client.Object.Copy(context.Background(), dest, sourceURL, nil)
	if err == nil {
		_, err = client.Object.Delete(context.Background(), source, nil)
		if err != nil {
			panic(err)
		}
	}
}
