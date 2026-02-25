package cas

import (
	"context"
	"errors"
	"fmt"

	cas "github.com/alibabacloud-go/cas-20200407/v4/client"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/model"
	"github.com/bububa/aliyun-acme-hook/internal/util"
)

func Upload(ctx context.Context, cfg *config.AliyunConfig, cert *model.Cert) error {
	if cert.FullChain == nil || cert.Key == nil {
		return errors.New("missing certificate path environment variables")
	}
	apiConfig, err := util.CreateOpenAPIClient(cfg.AK, cfg.SK, cfg.Region)
	if err != nil {
		return fmt.Errorf("failed to create CAS API config: %w", err)
	}
	client, err := cas.NewClient(apiConfig)
	if err != nil {
		return fmt.Errorf("failed to create CAS client: %w", err)
	}
	// 4. 构建上传请求
	uploadRequest := new(cas.UploadUserCertificateRequest)
	uploadRequest.SetName(cert.Name).SetCert(string(cert.FullChain)).SetKey(string(cert.Key))

	// 5. 执行上传
	result, err := client.UploadUserCertificate(uploadRequest)
	if err != nil {
		return fmt.Errorf("upload cert to CAS failed: %w", err)
	}
	cert.ID = *result.Body.CertId
	return nil
}
