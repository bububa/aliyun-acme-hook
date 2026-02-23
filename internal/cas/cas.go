package cas

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	cas "github.com/alibabacloud-go/cas-20200407/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/bububa/aliyun-acme-hook/config"
)

func Upload(ctx context.Context, cfg *config.AliyunConfig) (int64, error) {
	// 1. 获取 acme.sh 提供的环境变量
	// acme.sh 会在执行 hook 时自动 export 这些路径
	// certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("CERT_KEY_PATH")
	fullChainPath := os.Getenv("FULLCHAIN_PATH")
	domain := os.Getenv("CERT_DOMAIN") // acme.sh 传入的主域名

	if fullChainPath == "" || keyPath == "" {
		return 0, errors.New("missing certificate path environment variables")
	}

	// 2. 读取证书内容
	certContent, _ := os.ReadFile(fullChainPath)
	keyContent, _ := os.ReadFile(keyPath)

	// 3. 初始化阿里云客户端
	// 建议通过环境变量获取 Ali_Key 和 Ali_Secret
	config := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AK),
		AccessKeySecret: tea.String(cfg.SK),
		Endpoint:        tea.String(cfg.Region),
	}
	client, _ := cas.NewClient(config)

	// 4. 构建上传请求
	uploadRequest := &cas.UploadUserCertificateRequest{
		Name: tea.String(domain + "-" + time.Now().Format("20060102")), // 证书显示名称
		Cert: tea.String(string(certContent)),
		Key:  tea.String(string(keyContent)),
	}

	// 5. 执行上传
	result, err := client.UploadUserCertificate(uploadRequest)
	if err != nil {
		return 0, fmt.Errorf("upload cert to CAS failed: %w", err)
	}

	return *result.Body.CertId, nil
}
