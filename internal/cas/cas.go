package cas

import (
	"context"
	"errors"
	"fmt"

	cas "github.com/alibabacloud-go/cas-20200407/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/model"
)

func Upload(ctx context.Context, cfg *config.AliyunConfig, cert *model.Cert) error {
	if cert.FullChain == nil || cert.Key == nil {
		return errors.New("missing certificate path environment variables")
	}

	// 3. 初始化阿里云客户端
	// 建议通过环境变量获取 Ali_Key 和 Ali_Secret
	config := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AK),
		AccessKeySecret: tea.String(cfg.SK),
		RegionId:        tea.String(cfg.Region),
	}
	client, _ := cas.NewClient(config)
	// 4. 构建上传请求
	uploadRequest := &cas.UploadUserCertificateRequest{
		Name: tea.String(cert.Name), // 证书显示名称
		Cert: tea.String(string(cert.FullChain)),
		Key:  tea.String(string(cert.Key)),
	}

	// 5. 执行上传
	result, err := client.UploadUserCertificate(uploadRequest)
	if err != nil {
		return fmt.Errorf("upload cert to CAS failed: %w", err)
	}
	cert.ID = *result.Body.CertId
	return nil
}
