package cdn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	cdn "github.com/alibabacloud-go/cdn-20180510/v9/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/bububa/aliyun-acme-hook/config"
	"github.com/bububa/aliyun-acme-hook/internal/model"
	"github.com/bububa/aliyun-acme-hook/internal/util"
)

func Certificate(ctx context.Context, cfg *config.AliyunConfig, cert *model.Cert) error {
	if cert.Domain == "" {
		return errors.New("missing cert domain environment variables")
	}

	// 3. 初始化阿里云客户端
	// 建议通过环境变量获取 Ali_Key 和 Ali_Secret
	apiConfig, err := util.CreateOpenAPIClient(cfg.AK, cfg.SK, cfg.Region)
	if err != nil {
		return fmt.Errorf("failed to create CDN API config: %w", err)
	}
	client, err := cdn.NewClient(apiConfig)
	if err != nil {
		return fmt.Errorf("failed to create CDN client: %w", err)
	}
	domains, err := GetDomains(ctx, client, cert.Domain)
	if err != nil {
		return fmt.Errorf("certificate CDN domain failed, %w", err)
	}
	for _, domain := range domains {
		slog.InfoContext(ctx, "certicating CDN domain", "domain", domain)
		certReq := cdn.SetCdnDomainSSLCertificateRequest{
			DomainName:  tea.String(domain),
			CertName:    tea.String(cert.Name),
			CertId:      tea.Int64(cert.ID),
			CertType:    tea.String("cas"),
			SSLProtocol: tea.String("on"),
		}
		if _, err := client.SetCdnDomainSSLCertificate(&certReq); err != nil {
			slog.ErrorContext(ctx, "certicate CDN domain failed", "domain", domain, "error", err)
			return fmt.Errorf("set CDN domain cert failed, %w", err)
		}
		slog.InfoContext(ctx, "certicate CDN domain done!", "domain", domain)
	}
	return nil
}
